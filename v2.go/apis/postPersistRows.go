// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"fmt"

	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type persistGridRowDataFunc func(apiRequestParameters, *model.Grid, *model.Row) error

func persistGridRowData(r apiRequestParameters, grid *model.Grid, rows []*model.Row, f persistGridRowDataFunc) error {
	for _, row := range rows {
		err := f(r, grid, row)
		if err != nil {
			_ = r.rollbackTransaction()
			return err
		}
	}
	return nil
}

func postInsertGridRow(r apiRequestParameters, grid *model.Grid, row *model.Row) error {
	r.trace("postInsertGridRow()")
	row.TmpUuid = row.Uuid
	row.Uuid = utils.GetNewUUID()
	r.trace("postInsertGridRow() - row.TmpUuid=%v, row.Uuid=%v, row=%v", row.TmpUuid, row.Uuid, row)
	insertStatement := getInsertStatementForGridsApi(grid)
	insertValues := getInsertValuesForGridsApi(r.userUuid, grid, row)
	if err := r.execContext(insertStatement, insertValues...); err != nil {
		return r.logAndReturnError("Insert row error on %q: %v.", insertStatement, err)
	}
	r.log("Row [%s] inserted into %q.", row, grid.Uuid)
	return nil
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
		"revision, " +
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

func postUpdateGridRow(r apiRequestParameters, grid *model.Grid, row *model.Row) error {
	r.trace("postUpdateGridRow()")
	updateStatement := getUpdateStatementForGridsApi(grid)
	updateValues := getUpdateValuesForGridsApi(r.userUuid, grid, row)
	if err := r.execContext(updateStatement, updateValues...); err != nil {
		return r.logAndReturnError("Update row error on %q: %v.", updateStatement, err)
	}
	r.log("Row [%s] updated in %q.", row, grid.Uuid)
	return nil
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
		"revision = revision + 1, " +
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

func postDeleteGridRow(r apiRequestParameters, grid *model.Grid, row *model.Row) error {
	r.trace("postDeleteGridRow()")
	if err := r.execContext(getDeleteRowStatement(), row.Uuid, grid.Uuid); err != nil {
		return r.logAndReturnError("Delete row error: %v.", err)
	}
	r.log("Row [%s] deleted in %q.", row, grid.Uuid)
	return nil
}

func getDeleteRowStatement() string {
	return "DELETE FROM rows WHERE uuid = $1 and gridUuid = $2"
}
