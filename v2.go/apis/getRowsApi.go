// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func getGridsRows(ct context.Context, uri, dbName, gridUuid, uuid, userUuid, userName string) apiResponse {
	r, cancel, err := createContextAndApiRequestParameters(ct, dbName, userUuid, userName, uri)
	defer cancel()
	t := r.startTiming()
	defer r.stopTiming("getGridsRows()", t)
	if err != nil {
		return apiResponse{err: err}
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
		}
		canViewRows, canEditRows, canEditOwnedRows, canAddRows := grid.GetViewEditAccessFlags(r.userUuid)
		if !canViewRows {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Access forbidden."), forbidden: true}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, grid, uuid, true, true)
		if err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if uuid != "" && rowSetCount == 0 {
			r.ctxChan <- apiResponse{grid: grid, err: r.logAndReturnError("Data not found.")}
			return
		}
		r.ctxChan <- apiResponse{grid: grid, rows: rowSet, rowCount: rowSetCount, canViewRows: canViewRows, canEditRows: canEditRows, canEditOwnedRows: canEditOwnedRows, canAddRows: canAddRows}
		r.trace("getGridsRows() - Done")
	}()
	select {
	case <-r.ctx.Done():
		r.trace("getGridsRows() - Cancelled")
		return apiResponse{err: r.logAndReturnError("Get request has been cancelled: %v.", r.ctx.Err()), timeOut: true}
	case response := <-r.ctxChan:
		r.trace("getGridsRows() - OK")
		return response
	}
}

func getRowSetForGridsApi(r apiRequestParameters, grid *model.Grid, uuid string, getReferences bool, enabledOnly bool) ([]model.Row, int, error) {
	r.trace("getRowSetForGridsApi(%s, %s, %v)", uuid, grid, getReferences)
	t := r.startTiming()
	defer r.stopTiming("getRowSetForGridsApi()", t)
	query := getRowsQueryForGridsApi(grid, uuid, enabledOnly && uuid == "")
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
		r.trace("getRowSetForGridsApi(%s, %s, %v) - row=%v", uuid, grid, getReferences, row)
		gridForOwnership, err := getGridForOwnership(r, grid, row)
		if err != nil {
			return nil, 0, err
		}
		r.trace("getRowSetForGridsApi(%s, %s, %v) - gridForOwnership=%v", uuid, grid, getReferences, gridForOwnership)
		row.SetPathAndDisplayString(r.dbName)
		r.trace("getRowSetForGridsApi(%s, %s, %v) - row.DisplayString=%s", uuid, grid, getReferences, row.DisplayString)
		row.SetViewEditAccessFlags(gridForOwnership, r.userUuid)
		if getReferences {
			if err := getRelationshipsForRow(r, grid, row); err != nil {
				return nil, 0, err
			}
		}
		r.trace("getRowSetForGridsApi(%s, %s, %v) - row.CanViewRow=%v", uuid, grid, getReferences, row.CanViewRow)
		if row.CanViewRow {
			rows = append(rows, *row)
		}
	}
	return rows, len(rows), nil
}

// function is available for mocking
var getGridForOwnership = func(r apiRequestParameters, grid *model.Grid, row *model.Row) (*model.Grid, error) {
	r.trace("getGridForOwnership(%v, %v)", grid, row)
	if row.GridUuid == model.UuidGrids {
		return getGridForGridsApi(r, row.Uuid)
	} else if row.GridUuid == model.UuidRelationships && row.Text2 != nil {
		return getGridForGridsApi(r, *row.Text2)
	} else if row.GridUuid == model.UuidColumns {
		gridUuid, err := getGridUuidAttachedToColumn(r, row.Uuid)
		if err != nil {
			return nil, err
		}
		if gridUuid == "" {
			return nil, nil
		}
		return getGridForGridsApi(r, gridUuid)
	}
	return grid, nil
}

// function is available for mocking
var getRowsQueryForGridsApi = func(grid *model.Grid, uuid string, enabledOnly bool) string {
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
		"createdBy.text1, " +
		"rows.updated, " +
		"rows.updatedBy, " +
		"updatedBy.text1, " +
		"rows.revision " +

		"FROM " + grid.GetTableName() + " rows " +

		"LEFT OUTER JOIN users createdBy " +
		"ON createdBy.gridUuid = $1 " +
		"AND createdBy.uuid = rows.createdBy " +
		"AND createdBy.enabled = true " +

		"LEFT OUTER JOIN users updatedBy " +
		"ON updatedBy.gridUuid = $1 " +
		"AND updatedBy.uuid = rows.updatedBy " +
		"AND updatedBy.enabled = true " +

		getRowsWhereQueryForGridsApi(uuid) +
		getEnabledCondition(enabledOnly) +

		"ORDER BY rows.text1"
}

func getEnabledCondition(enabledOnly bool) string {
	if enabledOnly {
		return "AND rows.enabled = true "
	}
	return ""
}

func getRowsWhereQueryForGridsApi(uuid string) string {
	if uuid == "" {
		return "WHERE rows.griduuid = $2 "
	}
	return "WHERE rows.griduuid = $2 AND rows.uuid = $3 "
}

// function is available for mocking
var getRowsQueryParametersForGridsApi = func(gridUuid, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidUsers)
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
	output = append(output, &row.CreatedByName)
	output = append(output, &row.Updated)
	output = append(output, &row.UpdatedBy)
	output = append(output, &row.UpdatedByName)
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

func getGridUuidAttachedToColumn(r apiRequestParameters, uuid string) (string, error) {
	t := r.startTiming()
	defer r.stopTiming("getGridUuidAttachedToColumn()", t)
	var gridUuuid string
	query := getRowsQueryForGridUuidAttachedToColumn()
	parms := getRowsQueryParametersGridUuidAttachedToColumn(uuid)
	r.trace("getGridUuidAttachedToColumn(%s) - query=%s ; parms=%s", uuid, query, parms)
	if err := r.db.QueryRowContext(r.ctx, query, parms...).Scan(&gridUuuid); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", r.logAndReturnError("Error when retrieving grid uuid for column %q: %v.", uuid, err)
	}
	return gridUuuid, nil
}

// function is available for mocking
var getRowsQueryForGridUuidAttachedToColumn = func() string {
	return "SELECT text3 " +
		"FROM relationships " +
		"WHERE gridUuid = $1 " +
		"AND text2 = $2 " +
		"AND text4 = $3 " +
		"AND text5 = $4 " +
		"AND enabled = true"

}

func getRowsQueryParametersGridUuidAttachedToColumn(uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, model.UuidGrids)
	parameters = append(parameters, model.UuidColumns)
	parameters = append(parameters, uuid)
	return parameters
}
