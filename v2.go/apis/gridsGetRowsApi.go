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

func GetGridsRowsApi(c *gin.Context) {
	dbName := c.Param("dbName")
	gridUuid := c.Param("gridUuid")
	uuid := c.Param("uuid")
	_, user, err := getUserUuid(c, dbName)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	configuration.Trace(dbName, user, "GetGridsRowsApi()")
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUuid, uuid, user)
	if err != nil {
		c.Abort()
		if timeOut {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	configuration.Trace(dbName, user, "GetGridsRowsApi() - OK")
	if uuid != "" {
		c.JSON(http.StatusOK, gin.H{"uuid": uuid, "grid": grid, "rows": rowSet})
	} else {
		c.JSON(http.StatusOK, gin.H{"grid": grid, "rows": rowSet, "countRows": rowSetCount})
	}
}

func getUserUuid(c *gin.Context, dbName string) (string, string, error) {
	userUuid, user := c.GetString("userUuid"), c.GetString("user")
	auth, exists := c.Get("authorized")
	if len(userUuid) != len(model.UuidUsers) || user == "" || !exists || auth == false {
		return "", "", configuration.LogAndReturnError(dbName, user, "User not authorized.")
	}
	return userUuid, user, nil
}

func getGridsRows(ct context.Context, dbName, gridUuid, uuid, user string) (*model.Grid, []model.Row, int, bool, error) {
	configuration.Trace(dbName, user, "getGridsRows()")
	db, err := database.GetDbByName(dbName)
	if err != nil {
		return nil, nil, 0, false, err
	}
	ctxChan := make(chan apiResponse, 1)
	ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
	defer cancel()
	go func() {
		if err := database.TestSleep(ctx, dbName, user, db); err != nil {
			ctxChan <- apiResponse{err: configuration.LogAndReturnError(dbName, user, "Sleep interrupted: %v.", err)}
			return
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUuid)
		if err != nil {
			ctxChan <- apiResponse{err: err}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(ctx, db, dbName, user, uuid, grid, true)
		if uuid != "" && rowSetCount == 0 {
			ctxChan <- apiResponse{grid: grid, err: configuration.LogAndReturnError(dbName, user, "Data not found.")}
			return
		}
		ctxChan <- apiResponse{grid: grid, rows: rowSet, rowCount: rowSetCount}
	}()
	select {
	case <-ctx.Done():
		configuration.Trace(dbName, user, "getGridsRows() - Cancelled")
		return nil, nil, 0, true, configuration.LogAndReturnError(dbName, user, "Get request has been cancelled: %v.", ctx.Err())
	case response := <-ctxChan:
		configuration.Trace(dbName, user, "getGridsRows() - OK")
		return response.grid, response.rows, response.rowCount, false, response.err
	}
}

func getRowsQueryForGridsApi(grid *model.Grid, uuid string) string {
	selectStr := getRowsQueryColumnsForGridsApi(grid)
	fromStr := " FROM rows "
	whereStr := getRowsWhereQueryForGridsApi(uuid)
	orderByStr := " ORDER BY text1, text2, text3, text4 "
	return selectStr + fromStr + whereStr + orderByStr
}

func getRowsQueryColumnsForGridsApi(grid *model.Grid) string {
	var columns = ""
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			columns += col.Name + ", "
		}
	}
	return "SELECT uuid, " +
		"gridUuid, " +
		columns +
		"enabled, " +
		"created, " +
		"createdBy, " +
		"updated, " +
		"updatedBy, " +
		"revision "
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

func getRowSetForGridsApi(ctx context.Context, db *sql.DB, dbName, user, uuid string, grid *model.Grid, getReferences bool) ([]model.Row, int, error) {
	configuration.Trace(dbName, user, "getRowSetForGridsApi()")
	rows, err := db.QueryContext(ctx, getRowsQueryForGridsApi(grid, uuid), getRowsQueryParametersForGridsApi(grid.Uuid, uuid)...)
	if err != nil {
		return nil, 0, configuration.LogAndReturnError(dbName, user, "Error when querying rows: %v.", err)
	}
	defer rows.Close()
	var rowSet = make([]model.Row, 0)
	for rows.Next() {
		var row = new(model.Row)
		if err := rows.Scan(getRowsQueryOutputForGridsApi(grid, row)...); err != nil {
			return nil, 0, configuration.LogAndReturnError(dbName, user, "Error when scanning rows for %q: %v.", grid.Uuid, err)
		}
		row.SetPathAndDisplayString(dbName)
		if getReferences {
			if err := getRelationshipsForRow(ctx, db, dbName, user, grid, row); err != nil {
				return nil, 0, err
			}
		}
		rowSet = append(rowSet, *row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, configuration.LogAndReturnError(dbName, user, "Error when scanning rows for %q: %v.", grid.Uuid, err)
	}
	configuration.Trace(dbName, user, "Got %d rows from %q.", len(rowSet), grid.Uuid)
	return rowSet, len(rowSet), nil
}

func getRowsQueryOutputForGridsApi(grid *model.Grid, row *model.Row) []any {
	output := make([]any, 0)
	output = append(output, &row.Uuid)
	output = append(output, &row.GridUuid)
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			output = appendRowAttribute(output, row, col.Name)
		}
	}
	output = append(output, &row.Enabled)
	output = append(output, &row.Created)
	output = append(output, &row.CreatedBy)
	output = append(output, &row.Updated)
	output = append(output, &row.UpdatedBy)
	output = append(output, &row.Revision)
	return output
}

func appendRowAttribute(output []any, row *model.Row, attributeName string) []any {
	switch attributeName {
	case "text1":
		output = append(output, &row.Text1)
	case "text2":
		output = append(output, &row.Text2)
	case "text3":
		output = append(output, &row.Text3)
	case "text4":
		output = append(output, &row.Text4)
	case "text5":
		output = append(output, &row.Text5)
	case "text6":
		output = append(output, &row.Text6)
	case "text7":
		output = append(output, &row.Text7)
	case "text8":
		output = append(output, &row.Text8)
	case "text9":
		output = append(output, &row.Text9)
	case "text10":
		output = append(output, &row.Text10)
	case "int1":
		output = append(output, &row.Int1)
	case "int2":
		output = append(output, &row.Int2)
	case "int3":
		output = append(output, &row.Int3)
	case "int4":
		output = append(output, &row.Int4)
	case "int5":
		output = append(output, &row.Int5)
	case "int6":
		output = append(output, &row.Int6)
	case "int7":
		output = append(output, &row.Int7)
	case "int8":
		output = append(output, &row.Int8)
	case "int9":
		output = append(output, &row.Int9)
	case "int10":
		output = append(output, &row.Int10)
	}
	return output
}
