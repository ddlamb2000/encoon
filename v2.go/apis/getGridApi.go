// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"database/sql"

	"d.lambert.fr/encoon/model"
)

// function is available for mocking
var getGridForGridsApi = func(r apiRequestParameters, gridUuid string) (*model.Grid, error) {
	grid := new(model.Grid)
	query := getGridQueryForGridsApi()
	r.trace("getGridForGridsApi(%s) - query=%s", gridUuid, query)
	if err := r.db.QueryRowContext(r.ctx, query, model.UuidGrids, gridUuid).Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, r.logAndReturnError("Error when retrieving grid definition: %v.", err)
	}
	grid.SetPathAndDisplayString(r.dbName)
	err := getColumnsForGridsApi(r, grid)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

// function is available for mocking
var getGridQueryForGridsApi = func() string {
	return "SELECT uuid, " +
		"gridUuid, " +
		"text1, " +
		"text2, " +
		"text3, " +
		"text4, " +
		"enabled, " +
		"created, " +
		"createdBy, " +
		"updated, " +
		"updatedBy, " +
		"revision " +
		"FROM rows " +
		"WHERE gridUuid = $1 AND uuid = $2"
}

func getGridQueryOutputForGridsApi(grid *model.Grid) []any {
	output := make([]any, 0)
	output = append(output, &grid.Uuid)
	output = append(output, &grid.GridUuid)
	output = append(output, &grid.Text1)
	output = append(output, &grid.Text2)
	output = append(output, &grid.Text3)
	output = append(output, &grid.Text4)
	output = append(output, &grid.Enabled)
	output = append(output, &grid.Created)
	output = append(output, &grid.CreatedBy)
	output = append(output, &grid.Updated)
	output = append(output, &grid.UpdatedBy)
	output = append(output, &grid.Revision)
	return output
}

func getColumnsForGridsApi(r apiRequestParameters, grid *model.Grid) error {
	query := getGridColumsQueryForGridsApi()
	parms := getGridColumsQueryParametersForGridsApi(grid)
	r.trace("getColumnsForGridsApi(%s) - query=%s ; parms=%v", grid, query, parms)
	rows, err := r.queryContext(query, parms...)
	if err != nil {
		return r.logAndReturnError("Error when querying columns: %v.", err)
	}
	defer rows.Close()
	grid.Columns = make([]*model.Column, 0)
	for rows.Next() {
		var column = new(model.Column)
		if err := rows.Scan(getGridColumnQueryOutputForGridsApi(column)...); err != nil {
			return r.logAndReturnError("Error when scanning columns for: %v.", err)
		}
		r.trace("Got column for %s: %s.", grid, column)
		grid.Columns = append(grid.Columns, *&column)
	}
	return nil
}

// function is available for mocking
var getGridColumsQueryForGridsApi = func() string {
	return "SELECT col.text1, " +
		"col.text2, " +
		"coltype.uuid, " +
		"coltype.text1, " +
		"grid.uuid " +
		"FROM rows rel1 " +
		"INNER JOIN rows col " +
		"ON rel1.text4 = col.gridUuid " +
		"AND rel1.text5 = col.uuid " +
		"AND rel1.gridUuid = $1 " +
		"AND rel1.text1 = $2 " +
		"AND rel1.text2 = $3 " +
		"AND rel1.text3 = $4 " +
		"AND rel1.text4 = $5 " +
		"INNER JOIN rows rel2 " +
		"ON rel2.text2 = col.gridUuid " +
		"AND rel2.text3 = col.uuid " +
		"AND rel2.gridUuid = $6 " +
		"AND rel2.text1 = $7 " +
		"AND rel2.text4 = $8 " +
		"INNER JOIN rows coltype " +
		"ON rel2.text4 = coltype.gridUuid " +
		"AND rel2.text5 = coltype.Uuid " +
		"LEFT OUTER JOIN rows rel3 " +
		"ON rel3.text2 = col.gridUuid " +
		"AND rel3.text3 = col.uuid " +
		"AND rel3.gridUuid = $9 " +
		"AND rel3.text1 = $10 " +
		"LEFT OUTER JOIN rows grid " +
		"ON grid.gridUuid = rel3.text4 " +
		"AND grid.Uuid = rel3.text5 " +
		"ORDER BY col.text1 "
}

func getGridColumsQueryParametersForGridsApi(grid *model.Grid) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship1")
	parameters = append(parameters, model.UuidGrids)
	parameters = append(parameters, grid.Uuid)
	parameters = append(parameters, model.UuidColumns)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship1")
	parameters = append(parameters, model.UuidColumnTypes)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship2")
	return parameters
}

// function is available for mocking
var getGridColumnQueryOutputForGridsApi = func(column *model.Column) []any {
	output := make([]any, 0)
	output = append(output, &column.Label)
	output = append(output, &column.Name)
	output = append(output, &column.TypeUuid)
	output = append(output, &column.Type)
	output = append(output, &column.GridPromptUuid)
	return output
}
