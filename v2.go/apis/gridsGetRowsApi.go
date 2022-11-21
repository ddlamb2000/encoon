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
	configuration.Trace("GetGridsRowsApi()")
	_, user, err := getUserUui(c)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	dbName := c.Param("dbName")
	gridUri := c.Param("gridUri")
	uuid := c.Param("uuid")
	grid, rowSet, rowSetCount, timeOut, err := getGridsRows(c.Request.Context(), dbName, gridUri, uuid, user)
	if err != nil {
		c.Abort()
		if timeOut {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}
	configuration.Trace("GetGridsRowsApi() - OK")
	c.JSON(http.StatusOK, gin.H{"uuid": uuid, "grid": grid, "rows": rowSet, "countRows": rowSetCount})
}

func getUserUui(c *gin.Context) (string, string, error) {
	auth, exists := c.Get("authorized")
	if !exists || auth == false {
		return "", "", configuration.LogAndReturnError("Not authorized for %v.", c.Request.URL)
	}
	userUuid, user := c.GetString("userUuid"), c.GetString("user")
	if len(userUuid) < 10 || user == "" {
		return "", "", configuration.LogAndReturnError("User not authorized for %v.", c.Request.URL)
	}
	return userUuid, user, nil
}

type apiGetResponse struct {
	grid     *model.Grid
	rows     []model.Row
	rowCount int
	err      error
}

func getGridsRows(ct context.Context, dbName, gridUri, uuid, user string) (*model.Grid, []model.Row, int, bool, error) {
	configuration.Trace("getGridsRows()")
	db, err := database.GetDbByName(dbName)
	if err != nil {
		return nil, nil, 0, false, err
	}
	ctxChan := make(chan apiGetResponse, 1)
	ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
	defer cancel()
	go func() {
		if err := database.TestSleep(ctx, dbName, db); err != nil {
			ctxChan <- apiGetResponse{nil, nil, 0, configuration.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
			return
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, gridUri)
		if err != nil {
			ctxChan <- apiGetResponse{nil, nil, 0, err}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(ctx, db, dbName, user, uuid, grid, true)
		if uuid != "" && rowSetCount == 0 {
			ctxChan <- apiGetResponse{grid, rowSet, rowSetCount, configuration.LogAndReturnError("[%s] [%s] Data not found.", dbName, user)}
			return
		}
		ctxChan <- apiGetResponse{grid, rowSet, rowSetCount, err}
	}()
	select {
	case <-ctx.Done():
		configuration.Trace("getGridsRows() - Cancelled")
		return nil, nil, 0, true, configuration.LogAndReturnError("[%s] [%s] Get request has been cancelled: %v.", dbName, user, ctx.Err())
	case response := <-ctxChan:
		configuration.Trace("getGridsRows() - OK")
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

func getRowSetForGridsApi(ctx context.Context, db *sql.DB, dbName, user, uuid string, grid *model.Grid, getReferences bool) ([]model.Row, int, error) {
	configuration.Trace("getRowSetForGridsApi()")
	rows, err := db.QueryContext(ctx, getRowsQueryForGridsApi(grid, uuid), getRowsQueryParametersForGridsApi(grid.Uuid, uuid)...)
	if err != nil {
		return nil, 0, configuration.LogAndReturnError("[%s] [%s] Error when querying rows: %v.", dbName, user, err)
	}
	defer rows.Close()
	var rowSet = make([]model.Row, 0)
	for rows.Next() {
		var row = new(model.Row)
		if err := rows.Scan(getRowsQueryOutputForGridsApi(grid, row)...); err != nil {
			return nil, 0, configuration.LogAndReturnError("[%s] [%s] Error when scanning rows for %q: %v.", dbName, user, grid.GetUri(), err)
		}
		row.SetPathAndDisplayString(dbName, grid.GetUri())
		if getReferences {
			if err := getRelationshipsForRow(ctx, db, dbName, user, grid, row); err != nil {
				return nil, 0, err
			}
		}
		rowSet = append(rowSet, *row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, configuration.LogAndReturnError("[%s] [%s] Error when scanning rows for %q: %v.", dbName, user, grid.GetUri(), err)
	}
	configuration.Trace("[%s] [%s] Got %d rows from %q.", dbName, user, len(rowSet), grid.GetUri())
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
	output = append(output, &row.Version)
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
