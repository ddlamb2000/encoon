// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

func getGridForGridsApi(ctx context.Context, db *sql.DB, dbName, user, gridUriOrUuid, trace string) (*model.Grid, error) {
	grid := new(model.Grid)
	if err := db.QueryRowContext(ctx, getGridQueryForGridsApi(gridUriOrUuid), model.UuidGrids, gridUriOrUuid).
		Scan(getGridQueryOutputForGridsApi(grid)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.LogAndReturnError("[%s] [%s] Grid %q not found.", dbName, user, gridUriOrUuid)
		} else {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when retrieving grid definition %q: %v.", dbName, user, gridUriOrUuid, err)
		}
	}
	grid.SetPath(dbName, *grid.Text1)
	utils.Trace(trace, "Got grid %q: [%s].", gridUriOrUuid, grid)
	err := getColumnsForGridsApi(ctx, db, dbName, user, grid, trace)
	if err != nil {
		return nil, err
	}
	return grid, nil
}

func getGridQueryForGridsApi(gridUriOrUuid string) string {
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
		"version " +
		"FROM rows " +
		"WHERE gridUuid = $1 AND " +
		getGridQueryWhereForGridsApi(gridUriOrUuid)
}

func getGridQueryWhereForGridsApi(gridUriOrUuid string) string {
	if len(gridUriOrUuid) == len(model.UuidGrids) {
		return "uuid = $2"
	} else {
		return "text1 = $2"
	}
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
	output = append(output, &grid.Version)
	return output
}

func getColumnsForGridsApi(ctx context.Context, db *sql.DB, dbName, user string, grid *model.Grid, trace string) error {
	grid.Columns = make([]*model.Column, 0)
	statement := getGridColumsQueryForGridsApi()
	rows, err := db.QueryContext(ctx,
		statement,
		model.UuidRelationships,
		"relationship1",
		model.UuidGrids,
		grid.Uuid,
		model.UuidColumns,
		model.UuidRelationships,
		"relationship1",
		model.UuidColumnTypes)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Error when querying columns from %q using %q: %v.", dbName, user, grid.Uuid, statement, err)
	}
	defer rows.Close()
	for rows.Next() {
		var column = new(model.Column)
		if err := rows.Scan(getGridColumnQueryOutputForGridsApi(column)...); err != nil {
			return utils.LogAndReturnError("[%s] [%s] Error when scanning columns for %q using %q: %v.", dbName, user, grid.Uuid, statement, err)
		}
		utils.Trace(trace, "Got column for %q: [%s].", grid.Uuid, column)
		grid.Columns = append(grid.Columns, *&column)
	}
	return nil
}

func getGridColumsQueryForGridsApi() string {
	return "SELECT col.text1, " +
		"col.text2, " +
		"coltype.uuid, " +
		"coltype.text1 " +
		"FROM rows rel1 " +
		"INNER JOIN rows col " +
		"ON rel1.text4 = col.gridUuid " +
		"AND rel1.text5 = col.uuid " +
		"INNER JOIN rows rel2 " +
		"ON rel2.text2 = col.gridUuid " +
		"AND rel2.text3 = col.uuid " +
		"INNER JOIN rows coltype " +
		"ON rel2.text4 = coltype.gridUuid " +
		"AND rel2.text5 = coltype.Uuid " +
		"WHERE rel1.gridUuid = $1 " +
		"AND rel1.text1 = $2 " +
		"AND rel1.text2 = $3 " +
		"AND rel1.text3 = $4 " +
		"AND rel1.text4 = $5 " +
		"AND rel2.gridUuid = $6 " +
		"AND rel2.text1 = $7 " +
		"AND rel2.text4 = $8 " +
		"ORDER BY col.text1 "
}

func getGridColumnQueryOutputForGridsApi(column *model.Column) []any {
	output := make([]any, 0)
	output = append(output, &column.Label)
	output = append(output, &column.Name)
	output = append(output, &column.TypeUuid)
	output = append(output, &column.Type)
	return output
}
