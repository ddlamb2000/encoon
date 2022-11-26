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
	parms := getGridQueryParametersForGridsApi(gridUuid, r.userUuid)
	r.trace("getGridForGridsApi(%s) - query=%s ; parms=%v", gridUuid, query, parms)
	if err := r.db.QueryRowContext(r.ctx, query, parms...).Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, r.logAndReturnError("Error when retrieving grid definition: %v.", err)
	}
	grid.SetPathAndDisplayString(r.dbName)
	grid.SetViewEditAccessFlags(r.userUuid)
	err := getColumnsForGridsApi(r, grid)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

// function is available for mocking
var getGridQueryForGridsApi = func() string {
	return "SELECT grids.uuid, " +
		"grids.gridUuid, " +
		"grids.text1, " +
		"grids.text2, " +
		"grids.text3, " +
		"grids.enabled, " +
		"grids.created, " +
		"grids.createdBy, " +
		"grids.updated, " +
		"grids.updatedBy, " +
		"grids.revision, " +
		"owner.text5, " +
		"defAccess.text5, " +
		"viewAccess.text5, " +
		"editAccess.text5 " +
		"FROM grids " +

		"LEFT OUTER JOIN relationships owner " +
		"ON owner.gridUuid = $3 " +
		"AND owner.text1 = $4 " +
		"AND owner.text2 = $1 " +
		"AND owner.text3 = grids.uuid " +
		"AND owner.text4 = $5 " +
		"AND owner.text5 = $6 " +

		"LEFT OUTER JOIN relationships defAccess " +
		"ON defAccess.gridUuid = $3 " +
		"AND defAccess.text1 = $7 " +
		"AND defAccess.text2 = $1 " +
		"AND defAccess.text3 = grids.uuid " +
		"AND defAccess.text4 = $8 " +

		"LEFT OUTER JOIN relationships viewAccess " +
		"ON viewAccess.gridUuid = $3 " +
		"AND viewAccess.text1 = $9 " +
		"AND viewAccess.text2 = $1 " +
		"AND viewAccess.text3 = grids.uuid " +
		"AND viewAccess.text4 = $5 " +
		"AND viewAccess.text5 = $6 " +

		"LEFT OUTER JOIN relationships editAccess " +
		"ON editAccess.gridUuid = $3 " +
		"AND editAccess.text1 = $10 " +
		"AND editAccess.text2 = $1 " +
		"AND editAccess.text3 = grids.uuid " +
		"AND editAccess.text4 = $5 " +
		"AND editAccess.text5 = $6 " +

		"WHERE grids.gridUuid = $1 " +
		"AND grids.uuid = $2 "
}

func getGridQueryParametersForGridsApi(gridUuid, userUuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidGrids)
	parameters = append(parameters, gridUuid)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship3")
	parameters = append(parameters, model.UuidUsers)
	parameters = append(parameters, userUuid)
	parameters = append(parameters, "relationship2")
	parameters = append(parameters, model.UuidAccessLevel)
	parameters = append(parameters, "relationship4")
	parameters = append(parameters, "relationship5")
	return parameters
}

func getGridQueryOutputForGridsApi(grid *model.Grid) []any {
	output := make([]any, 0)
	output = append(output, &grid.Uuid)
	output = append(output, &grid.GridUuid)
	output = append(output, &grid.Text1)
	output = append(output, &grid.Text2)
	output = append(output, &grid.Text3)
	output = append(output, &grid.Enabled)
	output = append(output, &grid.Created)
	output = append(output, &grid.CreatedBy)
	output = append(output, &grid.Updated)
	output = append(output, &grid.UpdatedBy)
	output = append(output, &grid.Revision)
	output = append(output, &grid.OwnerUuid)
	output = append(output, &grid.DefaultAccessUuid)
	output = append(output, &grid.ViewAccessUuid)
	output = append(output, &grid.EditAccessUuid)
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
		"FROM relationships rel1 " +

		"INNER JOIN columns col " +
		"ON rel1.text4 = col.gridUuid " +
		"AND rel1.text5 = col.uuid " +
		"AND rel1.gridUuid = $1 " +
		"AND rel1.text1 = $2 " +
		"AND rel1.text2 = $3 " +
		"AND rel1.text3 = $4 " +
		"AND rel1.text4 = $5 " +

		"INNER JOIN relationships rel2 " +
		"ON rel2.text2 = col.gridUuid " +
		"AND rel2.text3 = col.uuid " +
		"AND rel2.gridUuid = $6 " +
		"AND rel2.text1 = $7 " +
		"AND rel2.text4 = $8 " +

		"INNER JOIN rows coltype " +
		"ON rel2.text4 = coltype.gridUuid " +
		"AND rel2.text5 = coltype.Uuid " +

		"LEFT OUTER JOIN relationships rel3 " +
		"ON rel3.text2 = col.gridUuid " +
		"AND rel3.text3 = col.uuid " +
		"AND rel3.gridUuid = $9 " +
		"AND rel3.text1 = $10 " +
		"LEFT OUTER JOIN grids grid " +
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
