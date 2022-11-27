// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func getGridsRows(ct context.Context, dbName, gridUuid, uuid, userUuid, userName string) (*model.Grid, []model.Row, int, apiResponse) {
	r, cancel, err := createContextAndApiRequestParameters(ct, dbName, userUuid, userName)
	defer cancel()
	if err != nil {
		return nil, nil, 0, apiResponse{err: err}
	}
	go func() {
		r.trace("getGridsRows()")
		database.Sleep(r.ctx, dbName, userName, r.db)
		grid, err := getGridForGridsApi(r, gridUuid)
		if err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		} else if grid == nil {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Data not found.")}
			return
		} else if !grid.CanViewRows {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Access forbidden."), forbidden: true}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, grid, uuid, true)
		if err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
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
		return nil, nil, 0, apiResponse{err: r.logAndReturnError("Get request has been cancelled: %v.", r.ctx.Err()), timeOut: true}
	case response := <-r.ctxChan:
		r.trace("getGridsRows() - OK")
		return response.grid, response.rows, response.rowCount, response
	}
}

func getRowSetForGridsApi(r apiRequestParameters, grid *model.Grid, uuid string, getReferences bool) ([]model.Row, int, error) {
	query := getRowsQueryForGridsApi(grid, uuid)
	parms := getRowsQueryParametersForGridsApi(grid.Uuid, uuid)
	r.trace("getRowSetForGridsApi(%s, %s, %v) - query=%s ; parms=%s", uuid, grid, getReferences, query, parms)
	set, err := r.db.QueryContext(r.ctx, query, parms...)
	if err != nil {
		return nil, 0, r.logAndReturnError("Error when querying rows: %v.", err)
	}
	defer set.Close()
	rows := make([]model.Row, 0)
	for set.Next() {
		row := model.GetNewRow()
		if err := set.Scan(getRowsQueryOutputForGridsApi(grid, row)...); err != nil {
			return nil, 0, r.logAndReturnError("Error when scanning rows: %v.", err)
		}
		gridForOwnership, err := getGridForOwnership(r, grid, row)
		if err != nil {
			return nil, 0, err
		}
		row.SetPathAndDisplayString(r.dbName)
		row.SetViewEditAccessFlags(gridForOwnership, r.userUuid)
		if getReferences {
			if err := getRelationshipsForRow(r, grid, row); err != nil {
				return nil, 0, err
			}
		}
		rows = append(rows, *row)
	}
	return rows, len(rows), nil
}

func getGridForOwnership(r apiRequestParameters, grid *model.Grid, row *model.Row) (*model.Grid, error) {
	if row.GridUuid == model.UuidGrids {
		return getGridForGridsApi(r, row.Uuid)
	} else if row.GridUuid == model.UuidRelationships && row.Text2 != nil {
		return getGridForGridsApi(r, *row.Text2)
	}
	return grid, nil
}

// function is available for mocking
var getRowsQueryForGridsApi = func(grid *model.Grid, uuid string) string {
	columns := ""
	for _, col := range grid.Columns {
		if col.IsAttribute() {
			columns += "rows." + col.Name + ", "
		}
	}
	return "SELECT rows.uuid, " +
		"rows.gridUuid, " +
		columns +
		"rows.enabled, " +
		"rows.created, " +
		"rows.createdBy, " +
		"rows.updated, " +
		"rows.updatedBy, " +
		"rows.revision " +
		"FROM " + grid.GetTableName() + " rows " +
		getRowsWhereQueryForGridsApi(uuid) +
		"ORDER BY rows.text1"
}

func getRowsWhereQueryForGridsApi(uuid string) string {
	if uuid == "" {
		return "WHERE rows.griduuid = $1 "
	}
	return "WHERE rows.griduuid = $1 AND rows.uuid = $2 "
}

// function is available for mocking
var getRowsQueryParametersForGridsApi = func(gridUuid, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, gridUuid)
	if uuid != "" {
		parameters = append(parameters, uuid)
	}
	return parameters
}

// function is available for mocking
var getRowsQueryOutputForGridsApi = func(grid *model.Grid, row *model.Row) []any {
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
