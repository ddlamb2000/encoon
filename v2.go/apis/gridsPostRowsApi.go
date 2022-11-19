// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

type gridPost struct {
	RowsAdded   []model.Row `json:"rowsAdded"`
	RowsEdited  []model.Row `json:"rowsEdited"`
	RowsDeleted []model.Row `json:"rowsDeleted"`
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
	var payload gridPost
	c.ShouldBindJSON(&payload)
	timeOut, err := postGridsRows(c.Request.Context(), dbName, userUuid, user, gridUri, payload.RowsAdded, payload.RowsEdited, payload.RowsDeleted, c.Query("trace"))
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
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUri, "", user, c.Query("trace"))
	if err != nil {
		c.Abort()
		if timeOut {
			utils.Trace(c.Query("trace"), "PostGridsRowsApi() - Timeout")
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			utils.Trace(c.Query("trace"), "PostGridsRowsApi() - Error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	utils.Trace(c.Query("trace"), "PostGridsRowsApi() - OK")
	c.JSON(http.StatusCreated, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

type apiPostResponse struct {
	err error
}

func postGridsRows(ct context.Context, dbName, userUuid, user, gridUri string, rowsAdded, rowsEdited, rowsDeleted []model.Row, trace string) (bool, error) {
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
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, rowsAdded, trace, postInsertGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, rowsEdited, trace, postUpdateGridRow); err != nil {
			ctxChan <- apiPostResponse{err}
			return
		}
		if err := persistGridRowData(ctx, dbName, db, userUuid, user, grid, rowsDeleted, trace, postDeleteGridRow); err != nil {
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

type persistGridRowDataFunc func(context.Context, string, *sql.DB, string, string, *model.Grid, *model.Row, string) error

func persistGridRowData(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, rows []model.Row, trace string, f persistGridRowDataFunc) error {
	for _, row := range rows {
		err := f(ctx, dbName, db, userUuid, user, grid, &row, trace)
		if err != nil {
			_ = RollbackTransaction(ctx, dbName, db, userUuid, user, trace)
			return err
		}
	}
	return nil
}

func postInsertGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, row *model.Row, trace string) error {
	utils.Trace(trace, "postInsertGridRow()")
	row.Uuid = utils.GetNewUUID()
	insertStatement := getInsertStatementForGridsApi(grid)
	insertValues := getInsertValuesForGridsApi(userUuid, grid, row)
	_, err := db.ExecContext(ctx, insertStatement, insertValues...)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Insert row error on %q: %v.", dbName, user, insertStatement, err)
	}
	utils.Log("[%s] [%s] Row [%s] inserted into %q.", dbName, user, row, grid.GetUri())
	return err
}

func getInsertStatementForGridsApi(grid *model.Grid) string {
	var parameterIndex = 4
	var columns, parameters = "", ""
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			columns += ", " + col.Name
			parameters += fmt.Sprintf(", $%d", parameterIndex)
			parameterIndex += 1
		}
	}
	insertStr := "INSERT INTO rows (uuid, " +
		"version, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid" +
		columns +
		") "
	valueStr := " VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3" +
		parameters +
		")"
	return insertStr + valueStr
}

func getInsertValuesForGridsApi(userUuid string, grid *model.Grid, row *model.Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	values = append(values, grid.Uuid)
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			values = appendRowParameter(values, row, col.Name)
		}
	}
	return values
}

func appendRowParameter(output []any, row *model.Row, attributeName string) []any {
	switch attributeName {
	case "text1":
		output = append(output, row.Text1)
	case "text2":
		output = append(output, row.Text2)
	case "text3":
		output = append(output, row.Text3)
	case "text4":
		output = append(output, row.Text4)
	case "text5":
		output = append(output, row.Text5)
	case "text6":
		output = append(output, row.Text6)
	case "text7":
		output = append(output, row.Text7)
	case "text8":
		output = append(output, row.Text8)
	case "text9":
		output = append(output, row.Text9)
	case "text10":
		output = append(output, row.Text10)
	case "int1":
		output = append(output, row.Int1)
	case "int2":
		output = append(output, row.Int2)
	case "int3":
		output = append(output, row.Int3)
	case "int4":
		output = append(output, row.Int4)
	case "int5":
		output = append(output, row.Int5)
	case "int6":
		output = append(output, row.Int6)
	case "int7":
		output = append(output, row.Int7)
	case "int8":
		output = append(output, row.Int8)
	case "int9":
		output = append(output, row.Int9)
	case "int10":
		output = append(output, row.Int10)
	}
	return output
}

func postUpdateGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, row *model.Row, trace string) error {
	utils.Trace(trace, "postUpdateGridRow()")
	updateStatement := getUpdateStatementForGridsApi(grid)
	updateValues := getUpdateValuesForGridsApi(userUuid, grid, row)
	_, err := db.ExecContext(ctx, updateStatement, updateValues...)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Update row error on %q: %v.", dbName, user, updateStatement, err)
	}
	utils.Log("[%s] [%s] Row [%s] updated in %q.", dbName, user, row, grid.GetUri())
	return err
}

func getUpdateStatementForGridsApi(grid *model.Grid) string {
	var parameterIndex = 4
	var columns = ""
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			columns += fmt.Sprintf(", %s = $%d", col.Name, parameterIndex)
			parameterIndex += 1
		}
	}
	updateStr := "UPDATE rows SET " +
		"version = version + 1, " +
		"updated = NOW(), " +
		"updatedBy = $3" +
		columns
	whereStr := " WHERE uuid = $1 and gridUuid = $2"
	return updateStr + whereStr
}

func getUpdateValuesForGridsApi(userUuid string, grid *model.Grid, row *model.Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, grid.Uuid)
	values = append(values, userUuid)
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			values = appendRowParameter(values, row, col.Name)
		}
	}
	return values
}

func postDeleteGridRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, row *model.Row, trace string) error {
	utils.Trace(trace, "postDeleteGridRow()")
	_, err := db.ExecContext(ctx, "DELETE FROM rows WHERE uuid = $1 and gridUuid = $2", row.Uuid, grid.Uuid)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Delete row error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Row [%s] deleted in %q.", dbName, user, row, grid.GetUri())
	return err
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
