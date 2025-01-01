// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"fmt"

	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type persistGridRowDataFunc func(ApiRequest, *model.Grid, *model.Row) error

func persistGridRowData(r ApiRequest, grid *model.Grid, rows []*model.Row, f persistGridRowDataFunc) error {
	for _, row := range rows {
		row.GridUuid = grid.Uuid
		err := f(r, grid, row)
		if err != nil {
			return err
		}
	}
	return nil
}

func postInsertGridRow(r ApiRequest, grid *model.Grid, row *model.Row) error {
	_, _, canAddRows, _ := grid.GetViewEditAccessFlags(r.p.UserUuid)
	if !canAddRows {
		r.ctxChan <- GridResponse{Err: r.logAndReturnError("Access forbidden."), Forbidden: true}
		return r.logAndReturnError("User isn't allowed to create rows.")
	}
	row.TmpUuid = row.Uuid
	if len(row.Uuid) != 36 {
		row.Uuid = utils.GetNewUUID()
	}
	query := getInsertStatementForGridsApi(grid)
	parms := getInsertValuesForGridsApi(r.p.UserUuid, grid, row)
	r.trace("postInsertGridRow(%s, %s) - query=%s, parms=%s", grid, row, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Insert row error: %v.", err)
	}
	r.log("Row [%s] inserted.", row.Uuid)
	return postInsertTransactionReferenceRow(r, grid, row, "relationship1")
}

// function is available for mocking
var getInsertStatementForGridsApi = func(grid *model.Grid) string {
	parameterIndex := 4
	columns, parameters := "", ""
	for _, col := range grid.Columns {
		if col.IsOwned() && col.IsAttribute() && col.Name != nil {
			columns += ", " + *col.Name
			parameters += fmt.Sprintf(", $%d", parameterIndex)
			parameterIndex += 1
		}
	}
	return "INSERT INTO " + grid.GetTableName() +
		" (uuid, " +
		"revision, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid" +
		columns +
		") " +
		"VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3" +
		parameters +
		")"
}

func getInsertValuesForGridsApi(userUuid string, grid *model.Grid, row *model.Row) []any {
	values := make([]any, 0)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	values = append(values, grid.Uuid)
	for _, col := range grid.Columns {
		if col.IsOwned() && col.IsAttribute() && col.Name != nil {
			values = appendRowParameter(values, row, *col.Name)
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

func postUpdateGridRow(r ApiRequest, grid *model.Grid, row *model.Row) error {
	rows, rowCount, err := getRowSetForGridsApi(r, grid, row.Uuid, false, true)
	if err != nil || rowCount != 1 {
		return r.logAndReturnError("Error retrieving row %q from grid %q before update: %v.", row.Uuid, grid.Uuid, err)
	}
	if !rows[0].CanEditRow {
		r.ctxChan <- GridResponse{Err: r.logAndReturnError("Access forbidden."), Forbidden: true}
		return r.logAndReturnError("User isn't allowed to update rows.")
	}
	query := getUpdateStatementForGridsApi(grid)
	parms := getUpdateValuesForGridsApi(r.p.UserUuid, grid, row)
	r.trace("postUpdateGridRow(%s, %s) - query=%s ; parms=%s", grid, row, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Update row error: %v.", err)
	}
	if err := removeAssociatedGridFromCache(r, grid, row.Uuid); err != nil {
		return r.logAndReturnError("Error when getting data for cache deletion: %v.", err)
	}
	r.log("Row [%s] updated.", row.Uuid)
	return postInsertTransactionReferenceRow(r, grid, row, "relationship2")
}

// function is available for mocking
var removeAssociatedGridFromCache = func(r ApiRequest, grid *model.Grid, uuid string) error {
	r.trace("removeAssociatedGridFromCache(%s, %v)", grid, uuid)
	if grid.Uuid == model.UuidGrids {
		r.trace("removeAssociatedGridFromCache() - Grid")
		removeGridFromCache(uuid)
		gridUuid, err := getGridUuidReferencedByColumn(r, uuid)
		if err != nil {
			return err
		}
		if gridUuid != "" {
			r.trace("removeAssociatedGridFromCache(%s, %v) - gridUuid=%s", grid, uuid, gridUuid)
			removeGridFromCache(gridUuid)
		}
	} else if grid.Uuid == model.UuidColumns {
		r.trace("removeAssociatedGridFromCache() - Column")
		gridUuid, err := getGridUuidAttachedToColumnForCache(r, uuid)
		if err != nil {
			return err
		}
		if gridUuid != "" {
			r.trace("removeAssociatedGridFromCache(%s, %v) - gridUuid=%s", grid, uuid, gridUuid)
			removeGridFromCache(gridUuid)
		}
	}
	return nil
}

// function is available for mocking
var removeAssociatedGridNotOwnedColumnFromCache = func(r ApiRequest, grid *model.Grid, uuid string) error {
	return removeAssociatedGridFromCache(r, grid, uuid)
}

// function is available for mocking
var getGridUuidAttachedToColumnForCache = func(r ApiRequest, uuid string) (string, error) {
	return getGridUuidAttachedToColumn(r, uuid)
}

// function is available for mocking
var getUpdateStatementForGridsApi = func(grid *model.Grid) string {
	parameterIndex := 4
	columns := ""
	for _, col := range grid.Columns {
		if col.IsOwned() && col.IsAttribute() && col.Name != nil {
			columns += fmt.Sprintf(", %s = $%d", *col.Name, parameterIndex)
			parameterIndex += 1
		}
	}
	return "UPDATE " +
		grid.GetTableName() +
		" SET revision = revision + 1, " +
		"updated = NOW(), " +
		"updatedBy = $3" +
		columns +
		" WHERE gridUuid = $1 " +
		"AND uuid = $2"
}

func getUpdateValuesForGridsApi(userUuid string, grid *model.Grid, row *model.Row) []any {
	values := make([]any, 0)
	values = append(values, grid.Uuid)
	values = append(values, row.Uuid)
	values = append(values, userUuid)
	for _, col := range grid.Columns {
		if col.IsOwned() && col.IsAttribute() && col.Name != nil {
			values = appendRowParameter(values, row, *col.Name)
		}
	}
	return values
}

func postDeleteGridRow(r ApiRequest, grid *model.Grid, row *model.Row) error {
	rows, rowCount, err := getRowSetForGridsApi(r, grid, row.Uuid, false, true)
	if err != nil || rowCount != 1 {
		return r.logAndReturnError("Error retrieving row %q from grid %q before delete: %v.", row.Uuid, grid.Uuid, err)
	}
	if !rows[0].CanEditRow {
		r.ctxChan <- GridResponse{Err: r.logAndReturnError("Access forbidden."), Forbidden: true}
		return r.logAndReturnError("User isn't allowed to update rows.")
	}
	query := getDeleteGridReferencedRowQuery(grid)
	r.trace("postDeleteGridRow(%s, %s) - query=%s", grid, row, query)
	if err := r.execContext(query, model.UuidRelationships, grid.Uuid, row.Uuid); err != nil {
		return r.logAndReturnError("Delete referenced row error: %v.", err)
	}

	query = getDeleteGridRowQuery(grid)
	r.trace("postDeleteGridRow(%s, %s) - query=%s", grid, row, query)
	if err := r.execContext(query, grid.Uuid, row.Uuid, r.p.UserUuid); err != nil {
		return r.logAndReturnError("Delete row error: %v.", err)
	}
	if err := removeAssociatedGridFromCache(r, grid, row.Uuid); err != nil {
		return r.logAndReturnError("Error when getting data for cache deletion: %v.", err)
	}
	r.log("Row [%s] deleted.", row.Uuid)
	return postInsertTransactionReferenceRow(r, grid, row, "relationship3")
}

// function is available for mocking
var getDeleteGridReferencedRowQuery = func(grid *model.Grid) string {
	return "UPDATE relationships set enabled = false WHERE gridUuid = $1 AND text2 = $2 AND text3 = $3"
}

// function is available for mocking
var getDeleteGridRowQuery = func(grid *model.Grid) string {
	return "UPDATE " +
		grid.GetTableName() +
		" SET revision = revision + 1, " +
		"enabled = false, " +
		"updated = NOW(), " +
		"updatedBy = $3 " +
		"WHERE gridUuid = $1 " +
		"AND uuid = $2"
}
