// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"d.lambert.fr/encoon/model"
)

// function is available for mocking
var getGridForGridsApi = func(r apiRequest, gridUuid string) (*model.Grid, error) {
	t := r.startTiming()
	defer r.stopTiming("getGridForGridsApi()", t)
	grid, ok := getGridFromCache(gridUuid)
	if ok && grid != nil {
		return grid, nil
	}
	grid, err := getGridInstanceForGridsApi(r, gridUuid)
	if err != nil || grid == nil {
		return nil, err
	}
	err = getColumnsForGridsApi(r, grid)
	if err != nil {
		return nil, err
	}
	cacheGrid(grid)
	return grid, nil
}

var getGridInstanceForGridsApi = func(r apiRequest, gridUuid string) (*model.Grid, error) {
	t := r.startTiming()
	defer r.stopTiming("getGridInstanceForGridsApi()", t)
	query := getGridQueryForGridsApi()
	parms := getGridQueryParametersForGridsApi(gridUuid, r.p.userUuid)
	r.trace("getGridInstanceForGridsApi(%s) - query=%s ; parms=%v", gridUuid, query, parms)
	set, err := r.db.QueryContext(r.ctx, query, parms...)
	if err != nil {
		return nil, r.logAndReturnError("Error when retrieving grid definition: %v.", err)
	}
	defer set.Close()
	grids := make([]model.Grid, 0)
	for set.Next() {
		grid := model.GetNewGrid()
		if err := set.Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
			return nil, r.logAndReturnError("Error when scanning grid definition: %v.", err)
		}
		grids = append(grids, *grid)
		grid.CopyAccessToOtherGrid(&grids[0])
	}
	if len(grids) == 0 {
		return nil, nil
	}
	grids[0].SetDisplayString(r.p.dbName)
	return &grids[0], nil
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
		"AND owner.enabled = true " +

		"LEFT OUTER JOIN relationships defAccess " +
		"ON defAccess.gridUuid = $3 " +
		"AND defAccess.text1 = $6 " +
		"AND defAccess.text2 = $1 " +
		"AND defAccess.text3 = grids.uuid " +
		"AND defAccess.text4 = $7 " +
		"AND defAccess.enabled = true " +

		"LEFT OUTER JOIN relationships viewAccess " +
		"ON viewAccess.gridUuid = $3 " +
		"AND viewAccess.text1 = $8 " +
		"AND viewAccess.text2 = $1 " +
		"AND viewAccess.text3 = grids.uuid " +
		"AND viewAccess.text4 = $5 " +
		"AND viewAccess.enabled = true " +

		"LEFT OUTER JOIN relationships editAccess " +
		"ON editAccess.gridUuid = $3 " +
		"AND editAccess.text1 = $9 " +
		"AND editAccess.text2 = $1 " +
		"AND editAccess.text3 = grids.uuid " +
		"AND editAccess.text4 = $5 " +
		"AND editAccess.enabled = true " +

		"WHERE grids.gridUuid = $1 " +
		"AND grids.uuid = $2 " +
		"AND grids.enabled = true"
}

func getGridQueryParametersForGridsApi(gridUuid, userUuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidGrids)
	parameters = append(parameters, gridUuid)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, "relationship3")
	parameters = append(parameters, model.UuidUsers)
	parameters = append(parameters, "relationship2")
	parameters = append(parameters, model.UuidAccessLevels)
	parameters = append(parameters, "relationship4")
	parameters = append(parameters, "relationship5")
	return parameters
}

// function is available for mocking
var getGridQueryOutputForGridsApi = func(grid *model.Grid) []any {
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

func getColumnsForGridsApi(r apiRequest, grid *model.Grid) error {
	t := r.startTiming()
	defer r.stopTiming("getColumnsForGridsApi()", t)
	grid.Columns = make([]*model.Column, 0)
	grid.Usages = make([]*model.Column, 0)
	queryOwned := getGridColumsOwnedQueryForGridsApi()
	parmsOwned := getGridColumsOwnedQueryParametersForGridsApi(grid)
	r.trace("getColumnsForGridsApi(%s) - queryOwned=%s ; parmsOwned=%v", grid, queryOwned, parmsOwned)
	if err := getColumnsRowsForGridsApi(r, grid, false, queryOwned, parmsOwned); err != nil {
		return err
	}
	queryNotOwned := getGridColumsNotOwnedQueryForGridsApi(true)
	parmsNotOwned := getGridColumsNotOwnedQueryParametersForGridsApi(grid)
	r.trace("getColumnsForGridsApi(%s) - queryNotOwned=%s ; parmsNotOwned=%v", grid, queryNotOwned, parmsNotOwned)
	if err := getColumnsRowsForGridsApi(r, grid, false, queryNotOwned, parmsNotOwned); err != nil {
		return err
	}
	queryUsages := getGridColumsNotOwnedQueryForGridsApi(false)
	r.trace("getColumnsForGridsApi(%s) - queryNotOwned=%s ; parmsNotOwned=%v", grid, queryUsages, parmsNotOwned)
	return getColumnsRowsForGridsApi(r, grid, true, queryUsages, parmsNotOwned)
}

func getColumnsRowsForGridsApi(r apiRequest, grid *model.Grid, setUsages bool, query string, parms []any) error {
	rows, err := r.queryContext(query, parms...)
	if err != nil {
		return r.logAndReturnError("Error when querying columns: %v.", err)
	}
	defer rows.Close()
	for rows.Next() {
		column := model.GetNewColumn()
		if err := rows.Scan(getGridColumnQueryOutputForGridsApi(column)...); err != nil {
			return r.logAndReturnError("Error when scanning columns for: %v.", err)
		}
		if !column.Owned {
			column.Grid, _ = getGridInstanceForGridsApi(r, column.GridUuid)
		}
		r.trace("Got column for %s: %s.", grid, column)
		if setUsages {
			grid.Usages = append(grid.Usages, column)
		} else {
			grid.Columns = append(grid.Columns, column)
		}
	}
	return nil
}

// function is available for mocking
var getGridColumsOwnedQueryForGridsApi = func() string {
	return "SELECT col.uuid, " +
		"col.text1, " +
		"col.text2, " +
		"coltype.uuid, " +
		"coltype.text1, " +
		"grid.uuid, " +
		"rel1.text3, " +
		"true " +
		"FROM relationships rel1 " +

		"INNER JOIN columns col " +
		"ON rel1.text4 = col.gridUuid " +
		"AND rel1.text5 = col.uuid " +
		"AND rel1.gridUuid = $1 " +
		"AND rel1.text1 = $2 " +
		"AND rel1.text2 = $3 " +
		"AND rel1.text3 = $4 " +
		"AND rel1.text4 = $5 " +
		"AND rel1.enabled = true " +

		"INNER JOIN relationships rel2 " +
		"ON rel2.text2 = col.gridUuid " +
		"AND rel2.text3 = col.uuid " +
		"AND rel2.gridUuid = $6 " +
		"AND rel2.text1 = $7 " +
		"AND rel2.text4 = $8 " +
		"AND rel2.enabled = true " +

		"INNER JOIN rows coltype " +
		"ON rel2.text4 = coltype.gridUuid " +
		"AND rel2.text5 = coltype.Uuid " +
		"AND coltype.enabled = true " +

		"LEFT OUTER JOIN relationships rel3 " +
		"ON rel3.text2 = col.gridUuid " +
		"AND rel3.text3 = col.uuid " +
		"AND rel3.gridUuid = $9 " +
		"AND rel3.text1 = $10 " +
		"AND rel3.enabled = true " +

		"LEFT OUTER JOIN grids grid " +
		"ON grid.gridUuid = rel3.text4 " +
		"AND grid.Uuid = rel3.text5 " +
		"AND grid.enabled = true " +

		"ORDER BY col.int1, col.text1 "
}

func getGridColumsOwnedQueryParametersForGridsApi(grid *model.Grid) []any {
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
var getGridColumsNotOwnedQueryForGridsApi = func(bidirectional bool) string {
	var bidirectionalCondition = ""
	if bidirectional {
		bidirectionalCondition = "AND col.text3 = 'true' "
	} else {
		bidirectionalCondition = "AND (col.text3 IS NULL OR col.text3 != 'true') "
	}
	return "SELECT col.uuid, " +
		"col.text1, " +
		"col.text2, " +
		"coltype.uuid, " +
		"coltype.text1, " +
		"rel1.text3, " +
		"rel1.text3, " +
		"false " +
		"FROM relationships rel1 " +

		"INNER JOIN columns col " +
		"ON rel1.text4 = col.gridUuid " +
		"AND rel1.text5 = col.uuid " +
		bidirectionalCondition +
		"AND rel1.gridUuid = $1 " +
		"AND rel1.text1 = $2 " +
		"AND rel1.text2 = $3 " +
		"AND rel1.text4 = $5 " +
		"AND rel1.enabled = true " +

		"INNER JOIN relationships rel2 " +
		"ON rel2.text2 = col.gridUuid " +
		"AND rel2.text3 = col.uuid " +
		"AND rel2.gridUuid = $6 " +
		"AND rel2.text1 = $7 " +
		"AND rel2.text4 = $8 " +
		"AND rel2.enabled = true " +

		"INNER JOIN rows coltype " +
		"ON rel2.text4 = coltype.gridUuid " +
		"AND rel2.text5 = coltype.Uuid " +
		"AND coltype.enabled = true " +

		"INNER JOIN relationships rel3 " +
		"ON rel3.text2 = col.gridUuid " +
		"AND rel3.text3 = col.uuid " +
		"AND rel3.gridUuid = $9 " +
		"AND rel3.text1 = $10 " +
		"AND rel3.text5 = $4 " +
		"AND rel3.enabled = true " +

		"ORDER BY col.int1, col.text1 "
}

func getGridColumsNotOwnedQueryParametersForGridsApi(grid *model.Grid) []any {
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
	output = append(output, &column.Uuid)
	output = append(output, &column.Label)
	output = append(output, &column.Name)
	output = append(output, &column.TypeUuid)
	output = append(output, &column.Type)
	output = append(output, &column.GridPromptUuid)
	output = append(output, &column.GridUuid)
	output = append(output, &column.Owned)
	return output
}
