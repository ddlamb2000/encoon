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
	"github.com/gin-gonic/gin"
)

type gridReferencePost struct {
	ColumnName string `json:"columnName"`
	FromUuid   string `json:"fromUuid"`
	ToGridUuid string `json:"toGridUuid"`
	ToUuid     string `json:"uuid"`
}

type gridPost struct {
	RowsAdded              []*model.Row        `json:"rowsAdded"`
	RowsEdited             []*model.Row        `json:"rowsEdited"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted"`
	ReferenceValuesAdded   []gridReferencePost `json:"referencedValuesAdded"`
	ReferenceValuesRemoved []gridReferencePost `json:"referencedValuesRemoved"`
}

func PostGridsRowsApi(c *gin.Context) {
	configuration.Trace("PostGridsRowsApi()")
	userUuid, user, err := getUserUui(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	var payload gridPost
	c.ShouldBindJSON(&payload)
	configuration.Trace("PostGridsRowsApi() - ReferenceValuesAdded=%v", payload.ReferenceValuesAdded)
	configuration.Trace("PostGridsRowsApi() - payload.RowsAdded=%v", payload.RowsAdded)
	timeOut, err := postGridsRows(c.Request.Context(), dbName, userUuid, user, gridUri, payload)
	if err != nil {
		c.Abort()
		if timeOut {
			configuration.Trace("PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			configuration.Trace("PostGridsRowsApi() - Error")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUri, "", user)
	if err != nil {
		c.Abort()
		if timeOut {
			configuration.Trace("PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			configuration.Trace("PostGridsRowsApi() - Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	configuration.Trace("PostGridsRowsApi() - OK")
	c.JSON(http.StatusCreated, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

type apiPostResponse struct {
	err error
}

func postGridsRows(ct context.Context, dbName, userUuid, user, gridUri string, payload gridPost) (bool, error) {
	configuration.Trace("postGridsRows()")
	db, err := database.GetDbByName(dbName)
	if err != nil {
		return false, err
	}
	ctxChan := make(chan apiPostResponse, 1)
	ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
	defer cancel()
	go func() {
		if err := database.TestSleep(ctx, dbName, db); err != nil {
			ctxChan <- apiPostResponse{configuration.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
			return
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri)
		if err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := BeginTransaction(ctx, dbName, db, userUuid, user); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, payload.RowsAdded, postInsertGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		configuration.Trace("postGridsRows() - payload.RowsAdded=%v", payload.RowsAdded)
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, payload.RowsEdited, postUpdateGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, payload.RowsDeleted, postDeleteGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridReferenceData(ctx, dbName, db, userUuid, user, grid, payload.RowsAdded, payload.ReferenceValuesAdded, postInsertReferenceRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridReferenceData(ctx, dbName, db, userUuid, user, grid, payload.RowsAdded, payload.ReferenceValuesRemoved, postDeleteReferenceRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := CommitTransaction(ctx, dbName, db, userUuid, user); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		ctxChan <- apiPostResponse{nil}
	}()
	configuration.Trace("postGridsRows() - Started")
	select {
	case <-ctx.Done():
		configuration.Trace("postGridsRows() - Cancelled")
		_ = RollbackTransaction(ctx, dbName, db, userUuid, user)
		return true, configuration.LogAndReturnError("[%s] [%s] Post request has been cancelled: %v.", dbName, user, ctx.Err())
	case response := <-ctxChan:
		configuration.Trace("postGridsRows() - OK ; response=%v", response)
		return false, response.err
	}
}

func BeginTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user string) error {
	configuration.Trace("beginTransaction()")
	_, err := db.ExecContext(ctx, "BEGIN")
	if err != nil {
		return configuration.LogAndReturnError("[%s] [%s] Begin transaction error: %v.", dbName, user, err)
	}
	configuration.Log("[%s] [%s] Begin transaction.", dbName, user)
	return err
}

func CommitTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user string) error {
	configuration.Trace("commitTransaction()")
	_, err := db.ExecContext(ctx, "COMMIT")
	if err != nil {
		return configuration.LogAndReturnError("[%s] [%s] Commit transaction error: %v.", dbName, user, err)
	}
	configuration.Log("[%s] [%s] Commit transaction.", dbName, user)
	return err
}

func RollbackTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user string) error {
	configuration.Trace("rollbackTransaction()")
	_, err := db.ExecContext(ctx, "ROLLBACK")
	if err != nil {
		return configuration.LogAndReturnError("[%s] [%s] Rollback transaction error: %v.", dbName, user, err)
	}
	configuration.Log("[%s] [%s] ROLLBACK transaction.", dbName, user)
	return err
}
