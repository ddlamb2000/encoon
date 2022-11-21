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
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	userUuid, user, err := getUserUui(c, dbName)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	configuration.Trace(dbName, user, "PostGridsRowsApi()")
	var payload gridPost
	c.ShouldBindJSON(&payload)
	configuration.Trace(dbName, user, "PostGridsRowsApi() - ReferenceValuesAdded=%v", payload.ReferenceValuesAdded)
	configuration.Trace(dbName, user, "PostGridsRowsApi() - payload.RowsAdded=%v", payload.RowsAdded)
	timeOut, err := postGridsRows(c.Request.Context(), dbName, userUuid, user, gridUri, payload)
	if err != nil {
		c.Abort()
		if timeOut {
			configuration.Trace(dbName, user, "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			configuration.Trace(dbName, user, "PostGridsRowsApi() - Error")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUri, "", user)
	if err != nil {
		c.Abort()
		if timeOut {
			configuration.Trace(dbName, user, "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			configuration.Trace(dbName, user, "PostGridsRowsApi() - Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	configuration.Trace(dbName, user, "PostGridsRowsApi() - OK")
	c.JSON(http.StatusCreated, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

type apiPostResponse struct {
	err error
}

func postGridsRows(ct context.Context, dbName, userUuid, user, gridUri string, payload gridPost) (bool, error) {
	configuration.Trace(dbName, user, "postGridsRows()")
	db, err := database.GetDbByName(dbName)
	if err != nil {
		return false, err
	}
	ctxChan := make(chan apiPostResponse, 1)
	ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
	defer cancel()
	go func() {
		if err := database.TestSleep(ctx, dbName, user, db); err != nil {
			ctxChan <- apiPostResponse{configuration.LogAndReturnError(dbName, user, "Sleep interrupted: %v.", err)}
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
		configuration.Trace(dbName, user, "postGridsRows() - payload.RowsAdded=%v", payload.RowsAdded)
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
	configuration.Trace(dbName, user, "postGridsRows() - Started")
	select {
	case <-ctx.Done():
		configuration.Trace(dbName, user, "postGridsRows() - Cancelled")
		_ = RollbackTransaction(ctx, dbName, db, userUuid, user)
		return true, configuration.LogAndReturnError(dbName, user, "Post request has been cancelled: %v.", ctx.Err())
	case response := <-ctxChan:
		configuration.Trace(dbName, user, "postGridsRows() - OK ; response=%v", response)
		return false, response.err
	}
}

func BeginTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user string) error {
	configuration.Trace(dbName, user, "beginTransaction()")
	_, err := db.ExecContext(ctx, "BEGIN")
	if err != nil {
		return configuration.LogAndReturnError(dbName, user, "Begin transaction error: %v.", err)
	}
	configuration.Log(dbName, user, "Begin transaction.")
	return err
}

func CommitTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user string) error {
	configuration.Trace(dbName, user, "commitTransaction()")
	_, err := db.ExecContext(ctx, "COMMIT")
	if err != nil {
		return configuration.LogAndReturnError(dbName, user, "Commit transaction error: %v.", err)
	}
	configuration.Log(dbName, user, "Commit transaction.")
	return err
}

func RollbackTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user string) error {
	configuration.Trace(dbName, user, "rollbackTransaction()")
	_, err := db.ExecContext(ctx, "ROLLBACK")
	if err != nil {
		return configuration.LogAndReturnError(dbName, user, "Rollback transaction error: %v.", err)
	}
	configuration.Log(dbName, user, "ROLLBACK transaction.")
	return err
}
