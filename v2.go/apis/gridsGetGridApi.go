// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

func getGridForGridsApi(ctx context.Context, db *sql.DB, dbName, user, gridUri, trace string) (*model.Grid, error) {
	grid := new(model.Grid)
	if err := db.QueryRowContext(ctx, getGridQueryForGridsApi(), model.UuidGrids, gridUri).
		Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogAndReturnError("[%s] [%s] Grid %q not found.", dbName, user, gridUri)
		} else {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when retrieving grid definition %q: %v.", dbName, user, gridUri, err)
		}
	}
	grid.SetPath(dbName, gridUri)
	utils.Trace(trace, "Got grid %q: [%s].", gridUri, grid)
	err := setColumnsForGridsApi(ctx, db, dbName, user, grid, trace)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

func getGridQueryForGridsApi() string {
	return "SELECT uuid, " +
		"text1, " +
		"text2, " +
		"text3, " +
		"text4, " +
		"enabled, " +
		"createdBy, " +
		"updatedBy, " +
		"version " +
		"FROM rows " +
		"WHERE gridUuid = $1 AND text1 = $2"
}

func getGridQueryOutputForGridsApi(grid *model.Grid) []any {
	output := make([]any, 0)
	output = append(output, &grid.Uuid)
	output = append(output, &grid.Text1)
	output = append(output, &grid.Text2)
	output = append(output, &grid.Text3)
	output = append(output, &grid.Text4)
	output = append(output, &grid.Enabled)
	output = append(output, &grid.CreateBy)
	output = append(output, &grid.UpdatedBy)
	output = append(output, &grid.Version)
	return output
}

func setColumnsForGridsApi(ctx context.Context, db *sql.DB, dbName, user string, grid *model.Grid, trace string) error {
	grid.Columns = make([]*model.Column, 0)
	statement := getGridColumsQueryForGridsApi()
	rows, err := db.QueryContext(ctx, statement, model.UuidRelationships, "relationship1", model.UuidGrids, grid.Uuid, model.UuidColumns)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Error when querying columns from %q using %q: %v.", dbName, user, grid.Uuid, statement, err)
	}
	defer rows.Close()
	for rows.Next() {
		var column = new(model.Column)
		if err := rows.Scan(getGridColumnQueryOutputForGridsApi(column)...); err != nil {
			return utils.LogAndReturnError("[%s] [%s] Error when scanning columns for %q using %q: %v.", dbName, user, grid.Uuid, statement, err)
		}
		utils.Trace(trace, "Got column for %q: [%s,%s,%s].", grid.Uuid, column.Text1, column.Text2, column.Text3)
		grid.Columns = append(grid.Columns, *&column)
	}
	return nil
}

func getGridColumsQueryForGridsApi() string {
	return "SELECT col.uuid, " +
		"col.gridUuid, " +
		"col.text1, " +
		"col.text2, " +
		"col.text3, " +
		"col.enabled, " +
		"col.createdBy, " +
		"col.updatedBy, " +
		"col.version " +
		"FROM rows rel " +
		"JOIN rows col " +
		"ON rel.text4 = col.gridUuid::text " +
		"AND rel.text5 = col.uuid::text " +
		"WHERE rel.gridUuid = $1 " +
		"AND rel.text1 = $2 " +
		"AND rel.text2 = $3 " +
		"AND rel.text3 = $4 " +
		"AND rel.text4 = $5 " +
		"ORDER BY col.text1 "
}

func getGridColumnQueryOutputForGridsApi(column *model.Column) []any {
	output := make([]any, 0)
	output = append(output, &column.Uuid)
	output = append(output, &column.GridUuid)
	output = append(output, &column.Text1)
	output = append(output, &column.Text2)
	output = append(output, &column.Text3)
	output = append(output, &column.Enabled)
	output = append(output, &column.CreateBy)
	output = append(output, &column.UpdatedBy)
	output = append(output, &column.Version)
	return output
}
