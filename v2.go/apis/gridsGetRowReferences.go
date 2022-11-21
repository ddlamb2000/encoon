// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

func getRelationshipsForRow(ctx context.Context, db *sql.DB, dbName, user string, grid *model.Grid, row *model.Row) error {
	configuration.Trace(dbName, user, "getRelationshipsForRow()")
	for _, col := range grid.Columns {
		if col.IsReference() {
			configuration.Trace(dbName, user, "getRelationshipsForRow() - col=%s", col)
			referencedRows, err := getReferencedRowsForRow(ctx, db, dbName, user, row, col.Name)
			if err != nil {
				return configuration.LogAndReturnError(dbName, user, "Error when retrieving referenced rows: %v.", err)
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

func getReferencedRowsForRow(ctx context.Context, db *sql.DB, dbName, user string, parentRow *model.Row, referenceName string) ([]model.Row, error) {
	statement := getQueryReferencedRowsForRow()
	parameters := getQueryParametersReferencedRowsForRow(referenceName, parentRow)
	rows, err := db.QueryContext(ctx, statement, parameters...)
	if err != nil {
		return nil, configuration.LogAndReturnError(dbName, user, "Error when querying referenced rows: %v.", err)
	}
	defer rows.Close()
	var rowSet = make([]model.Row, 0)
	for rows.Next() {
		var referencedUuid string
		var referencedGridUuid string
		if err := rows.Scan(&referencedGridUuid, &referencedUuid); err != nil {
			return nil, configuration.LogAndReturnError(dbName, user, "Error when scanning referenced rows: %v.", err)
		}
		grid, err := getGridForGridsApi(ctx, db, dbName, user, referencedGridUuid)
		if err != nil {
			return nil, configuration.LogAndReturnError(dbName, user, "Error when retrieving grid for referenced rows: %v.", err)
		}
		rows, _, err := getRowSetForGridsApi(ctx, db, dbName, user, referencedUuid, grid, false)
		if err != nil {
			return nil, configuration.LogAndReturnError(dbName, user, "Error when retrieving referenced rows: %v.", err)
		}
		rowSet = append(rowSet, rows...)
	}
	if err := rows.Err(); err != nil {
		return nil, configuration.LogAndReturnError(dbName, user, "Error when scanning referenced rows: %v.", err)
	}
	configuration.Trace(dbName, user, "Got referenced rows for %q.", parentRow)
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
