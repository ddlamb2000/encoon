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

func SeedDb(ct context.Context, dbName, importFileName string) error {
	db, err := GetDbByName(dbName)
	if err != nil {
		return err
	}
	return seedDb(ct, db, dbName, importFileName)
}

func seedDb(ctx context.Context, db *sql.DB, dbName string, seedDataFileName string) error {
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
	query := GetGridRowsQueryForSeedData(grid)
	parms := GetRowsQueryParametersForSeedData(row.GridUuid, row.Uuid)
	configuration.Trace(dbName, "", "seedRowDb(%s, %s) - query=%s, parms=%s", grid, row, query, parms)
	var uuid string
	var revision int64
	if err := db.QueryRowContext(ctx, query, parms...).Scan(&uuid, &revision); err != nil {
		if err == sql.ErrNoRows {
			configuration.Trace(dbName, "", "Not found row with grid uuid %q and uuid %q.", row.GridUuid, row.Uuid)
			return insertSeedRowDb(ctx, db, dbName, grid, &row)
		}
		return configuration.LogAndReturnError(dbName, "", "Error when retrieving row with grid uuid %q and uuid %q: %v.", row.GridUuid, row.Uuid, err)
	}
	if row.Uuid == uuid && row.Revision > revision {
		configuration.Trace(dbName, "", "seedRowDb(%s, %s) - found %q with revision %d that needs update.", grid, row, uuid, row.Revision)
		return updateSeedRowDb(ctx, db, dbName, grid, &row)
	}
	return nil
}

// function is available for mocking
var GetGridRowsQueryForSeedData = func(grid *model.Grid) string {
	return grid.GetRowsQueryForSeedData()
}

func GetRowsQueryParametersForSeedData(gridUuid, uuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, gridUuid)
	parameters = append(parameters, uuid)
	return parameters
}

func insertSeedRowDb(ctx context.Context, db *sql.DB, dbName string, grid *model.Grid, row *model.Row) error {
	query := GetGridInsertStatementForSeedRowDb(grid)
	parms := grid.GetInsertValuesForSeedRowDb(model.UuidRootUser, row)
	configuration.Trace(dbName, "", "insertSeedRowDb(%s, %s) - query=%s, parms=%s", grid, row, query, parms)
	if _, err := db.ExecContext(ctx, query, parms...); err != nil {
		return configuration.LogAndReturnError(dbName, "", "Insert error when seeding row: %v.", err)
	}
	configuration.Log(dbName, "", "Seed data row [%s] inserted.", row)
	return nil
}

// function is available for mocking
var GetGridInsertStatementForSeedRowDb = func(grid *model.Grid) string {
	return grid.GetInsertStatementForSeedRowDb()
}

func updateSeedRowDb(ctx context.Context, db *sql.DB, dbName string, grid *model.Grid, row *model.Row) error {
	query := GetGridUpdateStatementForSeedRowDb(grid)
	parms := grid.GetUpdateValuesForSeedRowDb(model.UuidRootUser, row)
	configuration.Trace(dbName, "", "updateSeedRowDb(%s, %s) - query=%s, parms=%s", grid, row, query, parms)
	if _, err := db.ExecContext(ctx, query, parms...); err != nil {
		return configuration.LogAndReturnError(dbName, "", "Update error when seeding row: %v.", err)
	}
	configuration.Log(dbName, "", "Seed data row [%s] with revision %d updated.", row, row.Revision)
	return nil
}

// function is available for mocking
var GetGridUpdateStatementForSeedRowDb = func(grid *model.Grid) string {
	return grid.GetUpdateStatementForSeedRowDb()
}
