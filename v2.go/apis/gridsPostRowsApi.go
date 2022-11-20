// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"
	"net/http"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

type gridReferencePost struct {
	Relationship string `json:"col"`
	FromUuid     string `json:"rowUuid"`
	ToGridUuid   string `json:"gridUuid"`
	ToUuid       string `json:"uuid"`
}

type gridPost struct {
	RowsAdded              []*model.Row        `json:"rowsAdded"`
	RowsEdited             []*model.Row        `json:"rowsEdited"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted"`
	ReferenceValuesAdded   []gridReferencePost `json:"referencedValuesAdded"`
	ReferenceValuesRemoved []gridReferencePost `json:"referencedValuesRemoved"`
}

func PostGridsRowsApi(c *gin.Context) {
	utils.Trace(c.Query("trace"), "PostGridsRowsApi()")
	userUuid, user, err := getUserUui(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	trace := c.Query("trace")
	var payload gridPost
	c.ShouldBindJSON(&payload)
	utils.Trace(trace, "PostGridsRowsApi() - ReferenceValuesAdded=%v", payload.ReferenceValuesAdded)
	utils.Trace(trace, "PostGridsRowsApi() - payload.RowsAdded=%v", payload.RowsAdded)
	timeOut, err := postGridsRows(c.Request.Context(), dbName, userUuid, user, gridUri, payload, trace)
	if err != nil {
		c.Abort()
		if timeOut {
			utils.Trace(trace, "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			utils.Trace(trace, "PostGridsRowsApi() - Error")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUri, "", user, trace)
	if err != nil {
		c.Abort()
		if timeOut {
			utils.Trace(trace, "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			utils.Trace(trace, "PostGridsRowsApi() - Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	utils.Trace(trace, "PostGridsRowsApi() - OK")
	c.JSON(http.StatusCreated, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

type apiPostResponse struct {
	err error
}

func postGridsRows(ct context.Context, dbName, userUuid, user, gridUri string, payload gridPost, trace string) (bool, error) {
	utils.Trace(trace, "postGridsRows()")
	db, err := database.GetDbByName(dbName)
	if err != nil {
		return false, err
	}
	ctxChan := make(chan apiPostResponse, 1)
	ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
	defer cancel()
	go func() {
		if err := database.TestSleep(ctx, dbName, db); err != nil {
			ctxChan <- apiPostResponse{utils.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
			return
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri, trace)
		if err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := BeginTransaction(ctx, dbName, db, userUuid, user, trace); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, payload.RowsAdded, trace, postInsertGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		utils.Trace(trace, "postGridsRows() - payload.RowsAdded=%v", payload.RowsAdded)
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, payload.RowsEdited, trace, postUpdateGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, payload.RowsDeleted, trace, postDeleteGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridReferenceData(ctx, dbName, db, userUuid, user, grid, payload.RowsAdded, payload.ReferenceValuesAdded, trace, postInsertReferenceRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridReferenceData(ctx, dbName, db, userUuid, user, grid, payload.RowsAdded, payload.ReferenceValuesRemoved, trace, postDeleteReferenceRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := CommitTransaction(ctx, dbName, db, userUuid, user, trace); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		ctxChan <- apiPostResponse{nil}
	}()
	utils.Trace(trace, "postGridsRows() - Started")
	select {
	case <-ctx.Done():
		utils.Trace(trace, "postGridsRows() - Cancelled")
		_ = RollbackTransaction(ctx, dbName, db, userUuid, user, trace)
		return true, utils.LogAndReturnError("[%s] [%s] Post request has been cancelled: %v.", dbName, user, ctx.Err())
	case response := <-ctxChan:
		utils.Trace(trace, "postGridsRows() - OK ; response=%v", response)
		return false, response.err
	}
}

func BeginTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user, trace string) error {
	utils.Trace(trace, "beginTransaction()")
	_, err := db.ExecContext(ctx, "BEGIN")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Begin transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Begin transaction.", dbName, user)
	return err
}

func CommitTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user, trace string) error {
	utils.Trace(trace, "commitTransaction()")
	_, err := db.ExecContext(ctx, "COMMIT")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Commit transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Commit transaction.", dbName, user)
	return err
}

func RollbackTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user, trace string) error {
	utils.Trace(trace, "rollbackTransaction()")
	_, err := db.ExecContext(ctx, "ROLLBACK")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Rollback transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] ROLLBACK transaction.", dbName, user)
	return err
}
