// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func getGridsRows(ct context.Context, dbName, gridUuid, uuid, userUuid, userName string) (*model.Grid, []model.Row, int, bool, error) {
	r, cancel, err := createContextAndApiRequestParameters(ct, dbName, userUuid, userName)
	defer cancel()
	if err != nil {
		return nil, nil, 0, false, err
	}
	go func() {
		r.trace("getGridsRows()")
		if err := database.TestSleep(r.ctx, dbName, userName, r.db); err != nil {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Sleep interrupted: %v.", err)}
			return
		}
		grid, err := getGridForGridsApi(r, gridUuid)
		if err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, uuid, grid, true)
		if uuid != "" && rowSetCount == 0 {
			r.ctxChan <- apiResponse{grid: grid, err: r.logAndReturnError("Data not found.")}
			return
		}
		r.ctxChan <- apiResponse{grid: grid, rows: rowSet, rowCount: rowSetCount}
		r.trace("getGridsRows() - Done")
	}()
	select {
	case <-r.ctx.Done():
		r.trace("getGridsRows() - Cancelled")
		return nil, nil, 0, true, r.logAndReturnError("Get request has been cancelled: %v.", r.ctx.Err())
	case response := <-r.ctxChan:
		r.trace("getGridsRows() - OK")
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

func getRowSetForGridsApi(r apiRequestParameters, uuid string, grid *model.Grid, getReferences bool) ([]model.Row, int, error) {
	r.trace("getRowSetForGridsApi()")
	rows, err := r.db.QueryContext(r.ctx, getRowsQueryForGridsApi(grid, uuid), getRowsQueryParametersForGridsApi(grid.Uuid, uuid)...)
	if err != nil {
		return nil, 0, r.logAndReturnError("Error when querying rows: %v.", err)
	}
	defer rows.Close()
	var rowSet = make([]model.Row, 0)
	for rows.Next() {
		var row = new(model.Row)
		if err := rows.Scan(getRowsQueryOutputForGridsApi(grid, row)...); err != nil {
			return nil, 0, r.logAndReturnError("Error when scanning rows for %q: %v.", grid.Uuid, err)
		}
		row.SetPathAndDisplayString(r.dbName)
		if getReferences {
			if err := getRelationshipsForRow(r, grid, row); err != nil {
				return nil, 0, err
			}
		}
		rowSet = append(rowSet, *row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, r.logAndReturnError("Error when scanning rows for %q: %v.", grid.Uuid, err)
	}
	r.trace("Got %d rows from %q.", len(rowSet), grid.Uuid)
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