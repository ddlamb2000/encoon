// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

func getRelationshipsForRow(ctx context.Context, db *sql.DB, dbName, user string, grid *model.Grid, row *model.Row, trace string) error {
	utils.Trace(trace, "getRelationshipsForRow()")
	for _, col := range grid.Columns {
		if col.IsReference() {
			utils.Trace(trace, "getRelationshipsForRow() - col=%s", col)
			referencedRows, err := getReferencedRowsForRow(ctx, db, dbName, user, row, col.Name, trace)
			if err != nil {
				return utils.LogAndReturnError("[%s] [%s] Error when retrieving referenced rows: %v.", dbName, user, err)
			}
			if len(referencedRows) > 0 {
				var reference = new(model.Reference)
				reference.Name = col.Name
				reference.Label = col.Label
				reference.Rows = referencedRows
				row.References = append(row.References, reference)
			}
		}
	}
	return nil
}

func getReferencedRowsForRow(ctx context.Context, db *sql.DB, dbName, user string, parentRow *model.Row, referenceName string, trace string) ([]model.Row, error) {
	statement := getQueryReferencedRowsForRow()
	parameters := getQueryParametersReferencedRowsForRow(referenceName, parentRow)
	rows, err := db.QueryContext(ctx, statement, parameters...)
	if err != nil {
		return nil, utils.LogAndReturnError("[%s] [%s] Error when querying referenced rows: %v.", dbName, user, err)
	}
	defer rows.Close()
	var rowSet = make([]model.Row, 0)
	for rows.Next() {
		var referencedUuid string
		var referencedGridUuid string
		if err := rows.Scan(&referencedGridUuid, &referencedUuid); err != nil {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when scanning referenced rows: %v.", dbName, user, err)
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, referencedGridUuid, trace)
		if err != nil {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when retrieving grid for referenced rows: %v.", dbName, user, err)
		}
		rows, _, err := getRowSetForGridsApi(ctx, db, dbName, user, referencedUuid, grid, false, trace)
		if err != nil {
			return nil, utils.LogAndReturnError("[%s] [%s] Error when retrieving referenced rows: %v.", dbName, user, err)
		}
		rowSet = append(rowSet, rows...)
	}
	if err := rows.Err(); err != nil {
		return nil, utils.LogAndReturnError("[%s] [%s] Error when scanning referenced rows: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Got referenced rows for %q.", dbName, user, parentRow)
	return rowSet, nil
}

func getQueryReferencedRowsForRow() string {
	return "SELECT text4, " +
		"text5 " +
		"FROM rows " +
		"WHERE gridUuid = $1 " +
		"AND text1 = $2 " +
		"AND text2 = $3 " +
		"AND text3 = $4"
}

func getQueryParametersReferencedRowsForRow(referenceName string, parentRow *model.Row) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, referenceName)
	parameters = append(parameters, parentRow.GridUuid)
	parameters = append(parameters, parentRow.Uuid)
	return parameters
}
