// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

const seedDataFileName = "../seedData/master.json"

func seedDb(ctx context.Context, db *sql.DB, dbName string) error {
	configuration.Log(dbName, "", "Start seeding data from %s.", seedDataFileName)
	f, err := os.ReadFile(seedDataFileName)
	if err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error loading seed data from file %q: %v.", seedDataFileName, err)
	}
	seedData := make([]model.Row, 0)
	if err = json.Unmarshal(f, &seedData); err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error parsing seed data from file %q: %v.", seedDataFileName, err)
	}
	for _, row := range seedData {
		row.SetDisplayString(dbName)
		if err := seedRowDb(ctx, db, dbName, row); err != nil {
			return err
		}
	}
	configuration.Log(dbName, "", "Data seeded.")
	return nil
}

func seedRowDb(ctx context.Context, db *sql.DB, dbName string, row model.Row) error {
	configuration.Trace(dbName, "", "Import %v.", row)
	grid := model.GetNewGrid(row.GridUuid)
	query := grid.GetRowsQueryForSeedData()
	parms := GetRowsQueryParametersForSeedData(row.GridUuid, row.Uuid)
	configuration.Trace(dbName, "", "seedRowDb(%s, %s) - query=%s, parms=%s", grid, row, query, parms)
	var uuid string
	if err := db.QueryRowContext(ctx, query, parms...).Scan(&uuid); err != nil {
		if err == sql.ErrNoRows {
			configuration.Trace(dbName, "", "Not found row with grid uuid %q and uuid %q.", row.GridUuid, row.Uuid)
			return insertSeedRowDb(ctx, db, dbName, grid, &row)
		}
		return configuration.LogAndReturnError(dbName, "", "Error when retrieving row with grid uuid %q and uuid %q: %v.", row.GridUuid, row.Uuid, err)
	}
	configuration.Trace(dbName, "", "Found %q.", uuid)
	return nil
}

func GetRowsQueryParametersForSeedData(gridUuid, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, gridUuid)
	parameters = append(parameters, uuid)
	return parameters
}

func insertSeedRowDb(ctx context.Context, db *sql.DB, dbName string, grid *model.Grid, row *model.Row) error {
	query := grid.GetInsertStatementForInsertSeedRowDb()
	parms := grid.GetInsertValuesForInsertSeedRowDb(model.UuidRootUser, row)
	configuration.Trace(dbName, "", "insertSeedRowDb(%s, %s) - query=%s, parms=%s", grid, row, query, parms)
	if _, err := db.ExecContext(ctx, query, parms...); err != nil {
		return configuration.LogAndReturnError(dbName, "", "Insert error when seeding row: %v.", err)
	}
	configuration.Log(dbName, "", "Row [%s] inserted.", row.Uuid)
	return nil
}
