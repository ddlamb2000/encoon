// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func GetGridsRows(ct context.Context, p ApiParameters, payload GridPost) GridResponse {
	r, cancel, err := createContextAndApiRequest(ct, p)
	defer cancel()
	t := r.startTiming()
	defer r.stopTiming("getGridsRows()", t)
	if err != nil {
		return GridResponse{Err: err}
	}
	go func() {
		r.trace("getGridsRows()")
		database.Sleep(r.ctx, p.DbName, p.UserName, r.db)
		grid, err := getGridForGridsApi(r, p.GridUuid)
		if err != nil {
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		} else if grid == nil {
			r.ctxChan <- GridResponse{Err: r.logAndReturnError("Data not found")}
			return
		}
		canViewRows, canEditRows, canAddRows, canEditGrid := grid.GetViewEditAccessFlags(p.UserUuid)
		if !canViewRows {
			r.ctxChan <- GridResponse{Err: r.logAndReturnError("Access forbidden"), Forbidden: true}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, grid, p.Uuid, true, true)
		if err != nil {
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if p.Uuid != "" && rowSetCount == 0 {
			r.ctxChan <- GridResponse{Grid: grid, Err: r.logAndReturnError("Data not found")}
			return
		}
		r.ctxChan <- GridResponse{
			Grid:                 grid,
			Rows:                 rowSet,
			CountRows:            rowSetCount,
			CanViewRows:          canViewRows,
			CanEditRows:          canEditRows,
			CanAddRows:           canAddRows,
			CanEditGrid:          canEditGrid,
			GridUuid:             p.GridUuid,
			ColumnUuid:           p.ColumnUuid,
			Uuid:                 p.Uuid,
			FilterColumnOwned:    p.FilterColumnOwned,
			FilterColumnName:     p.FilterColumnName,
			FilterColumnGridUuid: p.FilterColumnGridUuid,
			FilterColumnValue:    p.FilterColumnValue,
		}
		r.trace("getGridsRows() - Done")
	}()
	select {
	case <-r.ctx.Done():
		r.trace("getGridsRows() - Cancelled")
		return GridResponse{Err: r.logAndReturnError("Get request has been cancelled: %v.", r.ctx.Err()), TimeOut: true}
	case response := <-r.ctxChan:
		r.trace("getGridsRows() - OK")
		return response
	}
}

func getRowSetForGridsApi(r ApiRequest, grid *model.Grid, uuid string, getReferences bool, enabledOnly bool) ([]model.Row, int, error) {
	r.trace("getRowSetForGridsApi(%s, %s, %v)", grid, uuid, getReferences)
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
		r.trace("getRowSetForGridsApi(%s, %s, %v) - row=%v", grid, uuid, getReferences, row)
		gridForOwnership, err := getGridForOwnership(r, grid, row)
		if err != nil {
			return nil, 0, err
		}
		row.SetDisplayString(r.p.DbName)
		r.trace("getRowSetForGridsApi(%s, %s, %v) - row.DisplayString=%s", grid, uuid, getReferences, row.DisplayString)
		row.SetViewEditAccessFlags(gridForOwnership, r.p.UserUuid)
		if row.CanViewRow {
			if getReferences {
				references, err := getRelationshipsForRow(r, grid, row)
				if err != nil {
					return nil, 0, err
				}
				if uuid != "" {
					row.Audits, err = getAuditsForRow(r, grid, uuid)
					if err != nil {
						return nil, 0, err
					}
				}
				if matchesFilterColumn(r, references) {
					row.References = references
					rows = append(rows, *row)
				}
			} else {
				rows = append(rows, *row)
			}
		}
	}
	return rows, len(rows), nil
}

func matchesFilterColumn(r ApiRequest, references []*model.Reference) bool {
	if r.p.FilterColumnName == "" || r.p.FilterColumnGridUuid == "" || r.p.FilterColumnValue == "" {
		return true
	}
	for _, ref := range references {
		if ref.Owned == r.p.FilterColumnOwned && ref.Name == r.p.FilterColumnName {
			for _, refRow := range ref.Rows {
				r.trace("matchesFilterColumn() - refRow.GridUuid=%s ; r.p.filterColumnGridUuid=%s", refRow.GridUuid, r.p.FilterColumnGridUuid)
				if (ref.Owned || refRow.GridUuid == r.p.FilterColumnGridUuid) && refRow.Uuid == r.p.FilterColumnValue {
					return true
				}
			}
		}
	}
	return false
}

// function is available for mocking
var getGridForOwnership = func(r ApiRequest, grid *model.Grid, row *model.Row) (*model.Grid, error) {
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
	revisionEmbedding := ""
	if !grid.IsMetadata() {
		revisionEmbedding = ",rows.revisionEmbedding "
	}
	columns := ""
	for _, col := range grid.Columns {
		if col.IsOwned() && col.IsAttribute() && col.Name != nil && *col.Name != "" {
			columns += "rows." + *col.Name + ", "
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
		revisionEmbedding +
		"FROM " + grid.GetTableName() + " rows " +
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
		if col.IsOwned() && col.IsAttribute() && col.Name != nil && *col.Name != "" {
			output = appendRowAttribute(output, row, *col.Name)
		}
	}
	output = append(output, &row.Enabled)
	output = append(output, &row.Created)
	output = append(output, &row.CreatedBy)
	output = append(output, &row.Updated)
	output = append(output, &row.UpdatedBy)
	output = append(output, &row.Revision)
	if !grid.IsMetadata() {
		output = append(output, &row.RevisionEmbedding)
	}
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

// function is available for mocking
var getGridUuidAttachedToColumn = func(r ApiRequest, uuid string) (string, error) {
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
		"AND text1 = $2 " +
		"AND text2 = $3 " +
		"AND text4 = $4 " +
		"AND text5 = $5 " +
		"AND enabled = true"
}

func getRowsQueryParametersGridUuidAttachedToColumn(uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship1")
	parameters = append(parameters, model.UuidGrids)
	parameters = append(parameters, model.UuidColumns)
	parameters = append(parameters, uuid)
	return parameters
}

// function is available for mocking
var getGridUuidReferencedByColumn = func(r ApiRequest, uuid string) (string, error) {
	t := r.startTiming()
	defer r.stopTiming("getGridUuidReferencedByColumn()", t)
	var gridUuuid string
	query := getRowsQueryForGridUuidReferencedByColumn()
	parms := getRowsQueryParametersGridUuidReferencedByColumn(uuid)
	r.trace("getGridUuidReferencedByColumn(%s) - query=%s ; parms=%s", uuid, query, parms)
	if err := r.db.QueryRowContext(r.ctx, query, parms...).Scan(&gridUuuid); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", r.logAndReturnError("Error when retrieving grid uuid for column %q: %v.", uuid, err)
	}
	return gridUuuid, nil
}

// function is available for mocking
var getRowsQueryForGridUuidReferencedByColumn = func() string {
	return "SELECT text5 " +
		"FROM relationships " +
		"WHERE gridUuid = $1 " +
		"AND text1 = $2 " +
		"AND text2 = $3 " +
		"AND text3 = $4 " +
		"AND text4 = $5 " +
		"AND enabled = true"
}

func getRowsQueryParametersGridUuidReferencedByColumn(uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship2")
	parameters = append(parameters, model.UuidColumns)
	parameters = append(parameters, uuid)
	parameters = append(parameters, model.UuidGrids)
	return parameters
}
