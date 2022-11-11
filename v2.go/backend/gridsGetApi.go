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

func GetGridsRowsApi(c *gin.Context) {
	utils.Trace(c.Query("trace"), "GetGridsRowsApi()")
	_, user, err := getUserUui(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	logUri(c, dbName, user)
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(dbName, gridUri, uuid, user, c.Query("trace"))
	if err != nil {
		c.Abort()
		if timeOut {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	utils.Trace(c.Query("trace"), "GetGridsRowsApi() - OK")
	c.JSON(http.StatusOK, gin.H{"uuid": uuid, "grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

func getUserUui(c *gin.Context) (string, string, error) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		return "", "", utils.LogAndReturnError("Not authorized for %v.", c.Request.URL)
	}
	userUuid, user := c.GetString("userUuid"), c.GetString("user")
	if len(userUuid) < 10 || user == "" {
		return "", "", utils.LogAndReturnError("User not authorized for %v.", c.Request.URL)
	}
	return userUuid, user, nil
}

type apiGetResponse struct {
	grid     *Grid
	rows     []Row
	rowCount int
	err      error
}

func getGridsRows(dbName, gridUri, uuid, user, trace string) (*Grid, []Row, int, bool, error) {
	utils.Trace(trace, "getGridsRows()")
	db, err := getDbForGridsApi(dbName, user)
	if err != nil {
		return nil, nil, 0, false, err
	}
	ctxChan := make(chan apiGetResponse, 1)
	ctx, cancel := utils.GetContextWithTimeOut(dbName)
	defer cancel()
	go func() {
		if err := testSleep(ctx, dbName, db); err != nil {
			ctxChan <- apiGetResponse{nil, nil, 0, utils.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
			return
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri, trace)
		if err != nil {
			ctxChan <- apiGetResponse{nil, nil, 0, err}
			return
		}
		rows, err := getRowsForGridsApi(ctx, db, dbName, user, grid.Uuid, uuid, trace)
		if err != nil {
			ctxChan <- apiGetResponse{nil, nil, 0, err}
			return
		}
		defer rows.Close()
		rowSet, rowSetCount, err := getRowSetForGridsApi(dbName, user, gridUri, rows, trace)
		if uuid != "" && rowSetCount == 0 {
			ctxChan <- apiGetResponse{grid, rowSet, rowSetCount, utils.LogAndReturnError("[%s] [%s] Data not found.", dbName, user)}
			return
		}
		ctxChan <- apiGetResponse{grid, rowSet, rowSetCount, err}
	}()
	select {
	case <-ctx.Done():
		utils.Trace(trace, "getGridsRows() - Cancelled")
		return nil, nil, 0, true, utils.LogAndReturnError("[%s] [%s] Get request has been cancelled: %v.", dbName, user, ctx.Err())
	case response := <-ctxChan:
		utils.Trace(trace, "getGridsRows() - OK")
		return response.grid, response.rows, response.rowCount, false, response.err
	}
}

func getDbForGridsApi(dbName, user string) (*sql.DB, error) {
	if dbName == "" {
		return nil, utils.LogAndReturnError("[%s] Missing database name parameter.", user)
	}
	db := getDbByName(dbName)
	if db == nil {
		return nil, utils.LogAndReturnError("[%s] [%s] Database not available.", dbName, user)
	}
	return db, nil
}

func getGridForGridsApi(ctx context.Context, db *sql.DB, dbName, user, gridUri, trace string) (*Grid, error) {
	selectGridStatement := getGridQueryColumnsForGridsApi() + " FROM rows WHERE gridUuid = $1 AND text01 = $2"
	grid := new(Grid)
	if err := db.QueryRowContext(ctx, selectGridStatement, utils.UuidGrids, gridUri).
		Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogAndReturnError("[%s] [%s] Grid %q not found.", dbName, user, gridUri)
		} else {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when retrieving grid definition %q: %v.", dbName, user, gridUri, err)
		}
	}
	grid.SetPath(dbName, gridUri)
	utils.Trace(trace, "Got grid %q: [%s].", gridUri, grid)
	return grid, nil
}

func getGridQueryColumnsForGridsApi() string {
	return "SELECT uuid, " +
		"text01, " +
		"text02, " +
		"text03, " +
		"text04, " +
		"enabled, " +
		"createdBy, " +
		"updatedBy, " +
		"version"
}

func getGridQueryOutputForGridsApi(grid *Grid) []any {
	output := make([]any, 0)
	output = append(output, &grid.Uuid)
	output = append(output, &grid.Text01)
	output = append(output, &grid.Text02)
	output = append(output, &grid.Text03)
	output = append(output, &grid.Text04)
	output = append(output, &grid.Enabled)
	output = append(output, &grid.CreateBy)
	output = append(output, &grid.UpdatedBy)
	output = append(output, &grid.Version)
	return output
}

func getRowsForGridsApi(ctx context.Context, db *sql.DB, dbName, user, gridUuid, uuid, trace string) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx,
		getRowsQueryForGridsApi(uuid),
		getRowsQueryParametersForGridsApi(gridUuid, uuid)...)
	if err != nil {
		return nil, utils.LogAndReturnError("[%s] [%s] Error when querying rows: %v.", dbName, user, err)
	}
	return rows, nil
}

func getRowsQueryForGridsApi(uuid string) string {
	selectStr := getRowsQueryColumnsForGridsApi()
	fromStr := " FROM rows "
	whereStr := getRowsWhereQueryForGridsApi(uuid)
	orderByStr := " ORDER BY text01, text02, text03, text04 "
	return selectStr + fromStr + whereStr + orderByStr
}

func getRowsQueryColumnsForGridsApi() string {
	return "SELECT uuid, " +
		"text01, " +
		"text02, " +
		"text03, " +
		"text04, " +
		"int01, " +
		"int02, " +
		"int03, " +
		"int04, " +
		"enabled, " +
		"createdBy, " +
		"updatedBy, " +
		"version "
}

func getRowsWhereQueryForGridsApi(uuid string) string {
	if uuid == "" {
		return " WHERE griduuid = $1 "
	}
	return " WHERE uuid = $2 AND griduuid = $1 "
}

func getRowsQueryParametersForGridsApi(gridUuid, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, gridUuid)
	if uuid != "" {
		parameters = append(parameters, uuid)
	}
	return parameters
}

func getRowSetForGridsApi(dbName, user, gridUri string, rows *sql.Rows, trace string) ([]Row, int, error) {
	var rowSet = make([]Row, 0)
	var rowSetCount = 0
	for rows.Next() {
		var row = new(Row)
		if err := rows.Scan(getRowsQueryOutputForGridsApi(row)...); err != nil {
			return nil, 0, utils.LogAndReturnError("[%s] [%s] Error when scanning rows for %q: %v.", dbName, user, gridUri, err)
		}
		row.SetPath(dbName, gridUri)
		rowSet = append(rowSet, *row)
		rowSetCount += 1
	}
	if err := rows.Err(); err != nil {
		return nil, 0, utils.LogAndReturnError("[%s] [%s] Error when scanning rows for %q: %v.", dbName, user, gridUri, err)
	}
	utils.Log("[%s] [%s] Got %d rows from %q.", dbName, user, rowSetCount, gridUri)
	return rowSet, rowSetCount, nil
}

func getRowsQueryOutputForGridsApi(row *Row) []any {
	output := make([]any, 0)
	output = append(output, &row.Uuid)
	output = append(output, &row.Text01)
	output = append(output, &row.Text02)
	output = append(output, &row.Text03)
	output = append(output, &row.Text04)
	output = append(output, &row.Int01)
	output = append(output, &row.Int02)
	output = append(output, &row.Int03)
	output = append(output, &row.Int04)
	output = append(output, &row.Enabled)
	output = append(output, &row.CreateBy)
	output = append(output, &row.UpdatedBy)
	output = append(output, &row.Version)
	return output
}
