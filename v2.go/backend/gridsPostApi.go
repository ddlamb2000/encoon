// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"database/sql"
	"net/http"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

type gridPost struct {
	RowsAdded   []Row `json:"rowsAdded"`
	RowsEdited  []Row `json:"rowsEdited"`
	RowsDeleted []Row `json:"rowsDeleted"`
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
	logUri(c, dbName, user)
	var payload gridPost
	c.ShouldBindJSON(&payload)
	timeOut, err := postGridsRows(dbName, userUuid, user, gridUri, payload.RowsAdded, payload.RowsEdited, payload.RowsDeleted, c.Query("trace"))
	if err != nil {
		c.Abort()
		if timeOut {
			utils.Trace(c.Query("trace"), "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			utils.Trace(c.Query("trace"), "PostGridsRowsApi() - Error")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(dbName, gridUri, "", user, c.Query("trace"))
	if err != nil {
		c.Abort()
		if timeOut {
			utils.Trace(c.Query("trace"), "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			utils.Trace(c.Query("trace"), "PostGridsRowsApi() - Error")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	utils.Trace(c.Query("trace"), "PostGridsRowsApi() - OK")
	c.JSON(http.StatusOK, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

type apiPostResponse struct {
	err error
}

func postGridsRows(dbName, userUuid, user, gridUri string, rowsAdded []Row, rowsEdited []Row, rowsDeleted []Row, trace string) (bool, error) {
	utils.Trace(trace, "postGridsRows()")
	db, err := getDbForGridsApi(dbName, user)
	if err != nil {
		return false, err
	}
	ctxChan := make(chan apiPostResponse, 1)
	ctx, cancel := utils.GetContextWithTimeOut(dbName)
	defer cancel()
	go func() {
		if err := testSleep(ctx, dbName, db); err != nil {
			ctxChan <- apiPostResponse{utils.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
			return
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri, trace)
		if err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := beginTransaction(ctx, dbName, db, userUuid, user, trace); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		for _, row := range rowsAdded {
			err := postInsertGridRow(ctx, dbName, db, userUuid, user, gridUri, grid.Uuid, row, trace)
			if err != nil {
				_ = rollbackTransaction(ctx, dbName, db, userUuid, user, trace)
				ctxChan <- apiPostResponse{err}
				return
			}
		}
		for _, row := range rowsEdited {
			err := postUpdateGridRow(ctx, dbName, db, userUuid, user, gridUri, grid.Uuid, row, trace)
			if err != nil {
				_ = rollbackTransaction(ctx, dbName, db, userUuid, user, trace)
				ctxChan <- apiPostResponse{err}
				return
			}
		}
		for _, row := range rowsDeleted {
			err := postDeleteGridRow(ctx, dbName, db, userUuid, user, gridUri, grid.Uuid, row, trace)
			if err != nil {
				_ = rollbackTransaction(ctx, dbName, db, userUuid, user, trace)
				ctxChan <- apiPostResponse{err}
				return
			}
		}
		if err := commitTransaction(ctx, dbName, db, userUuid, user, trace); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		ctxChan <- apiPostResponse{nil}
	}()
	utils.Trace(trace, "postGridsRows() - Started")
	select {
	case <-ctx.Done():
		utils.Trace(trace, "postGridsRows() - Cancelled")
		_ = rollbackTransaction(ctx, dbName, db, userUuid, user, trace)
		return true, utils.LogAndReturnError("[%s] [%s] Post request has been cancelled: %v.", dbName, user, ctx.Err())
	case response := <-ctxChan:
		utils.Trace(trace, "postGridsRows() - OK ; response=%v", response)
		return false, response.err
	}
}

func postInsertGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user, gridUri, gridUuid string, row Row, trace string) error {
	utils.Trace(trace, "postInsertGridRow()")
	row.Uuid = utils.GetNewUUID()
	insertStatement := getInsertStatementForGridsApi()
	insertValues := getInsertValuesForGridsApi(userUuid, gridUuid, row)
	_, err := db.ExecContext(ctx, insertStatement, insertValues...)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Insert row error on %q: %v.", dbName, user, insertStatement, err)
	}
	utils.Log("[%s] [%s] Row [%s] inserted into %q.", dbName, user, row, gridUri)
	return err
}

func getInsertStatementForGridsApi() string {
	insertStr := "INSERT INTO rows (uuid, " +
		"version, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"text01, " +
		"text02, " +
		"text03, " +
		"text04, " +
		"int01, " +
		"int02, " +
		"int03, " +
		"int04) "
	valueStr := " VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3, " +
		"$4, " +
		"$5, " +
		"$6, " +
		"$7, " +
		"$8, " +
		"$9, " +
		"$10, " +
		"$11)"
	return insertStr + valueStr
}

func getInsertValuesForGridsApi(userUuid, gridUuid string, row Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	values = append(values, gridUuid)
	values = append(values, row.Text01)
	values = append(values, row.Text02)
	values = append(values, row.Text03)
	values = append(values, row.Text04)
	values = append(values, row.Int01)
	values = append(values, row.Int02)
	values = append(values, row.Int03)
	values = append(values, row.Int04)
	return values
}

func postUpdateGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user, gridUri, gridUuid string, row Row, trace string) error {
	utils.Trace(trace, "postUpdateGridRow()")
	updateStatement := getUpdateStatementForGridsApi()
	updateValues := getUpdateValuesForGridsApi(userUuid, gridUuid, row)
	_, err := db.ExecContext(ctx, updateStatement, updateValues...)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Update row error on %q: %v.", dbName, user, updateStatement, err)
	}
	utils.Log("[%s] [%s] Row [%s] updated in %q.", dbName, user, row, gridUri)
	return err
}

func getUpdateStatementForGridsApi() string {
	updateStr := "UPDATE rows SET " +
		"version = version + 1, " +
		"updated = NOW(), " +
		"updatedBy = $3, " +
		"text01 = $4, " +
		"text02 = $5, " +
		"text03 = $6, " +
		"text04 = $7, " +
		"int01 = $8, " +
		"int02 = $9, " +
		"int03 = $10, " +
		"int04 = $11"
	whereStr := " WHERE uuid = $1 and gridUuid = $2"
	return updateStr + whereStr
}

func getUpdateValuesForGridsApi(userUuid, gridUuid string, row Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, gridUuid)
	values = append(values, userUuid)
	values = append(values, row.Text01)
	values = append(values, row.Text02)
	values = append(values, row.Text03)
	values = append(values, row.Text04)
	values = append(values, row.Int01)
	values = append(values, row.Int02)
	values = append(values, row.Int03)
	values = append(values, row.Int04)
	return values
}

func postDeleteGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user, gridUri, gridUuid string, row Row, trace string) error {
	utils.Trace(trace, "postDeleteGridRow()")
	_, err := db.ExecContext(ctx, "DELETE FROM rows WHERE uuid = $1 and gridUuid = $2", row.Uuid, gridUuid)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Delete row error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Row [%s] deleted in %q.", dbName, user, row, gridUri)
	return err
}
