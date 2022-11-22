// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

func getGridForGridsApi(ctx context.Context, db *sql.DB, dbName, user, gridUuid string) (*model.Grid, error) {
	grid := new(model.Grid)
	if err := db.QueryRowContext(ctx, getGridQueryForGridsApi(), model.UuidGrids, gridUuid).
		Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, configuration.LogAndReturnError(dbName, user, "Grid %q not found.", gridUuid)
		} else {
			return nil, configuration.LogAndReturnError(dbName, user, "Error when retrieving grid definition %q: %v.", gridUuid, err)
		}
	}
	grid.SetPath(dbName)
	configuration.Trace(dbName, user, "Got grid %q: [%s].", gridUuid, grid)
	err := getColumnsForGridsApi(ctx, db, dbName, user, grid)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

func getGridQueryForGridsApi() string {
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

func getColumnsForGridsApi(ctx context.Context, db *sql.DB, dbName, user string, grid *model.Grid) error {
	grid.Columns = make([]*model.Column, 0)
	statement := getGridColumsQueryForGridsApi()
	configuration.Trace(dbName, user, "getColumnsForGridsApi() - statement=%s", statement)
	rows, err := db.QueryContext(
		ctx,
		statement,
		model.UuidRelationships,
		"relationship1",
		model.UuidGrids,
		grid.Uuid,
		model.UuidColumns,
		model.UuidRelationships,
		"relationship1",
		model.UuidColumnTypes,
		model.UuidRelationships,
		"relationship2",
	)
	if err != nil {
		return configuration.LogAndReturnError(dbName, user, "Error when querying columns from %q using %q: %v.", grid.Uuid, statement, err)
	}
	defer rows.Close()
	for rows.Next() {
		var column = new(model.Column)
		if err := rows.Scan(getGridColumnQueryOutputForGridsApi(column)...); err != nil {
			return configuration.LogAndReturnError(dbName, user, "Error when scanning columns for %q using %q: %v.", grid.Uuid, statement, err)
		}
		configuration.Trace(dbName, user, "Got column for %q: [%s].", grid.Uuid, column)
		grid.Columns = append(grid.Columns, *&column)
	}
	return nil
}

func getGridColumsQueryForGridsApi() string {
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

func getGridColumnQueryOutputForGridsApi(column *model.Column) []any {
	output := make([]any, 0)
	output = append(output, &column.Label)
	output = append(output, &column.Name)
	output = append(output, &column.TypeUuid)
	output = append(output, &column.Type)
	output = append(output, &column.GridPromptUuid)
	return output
}
