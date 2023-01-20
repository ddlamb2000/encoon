// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

func migrateDb(ctx context.Context, db *sql.DB, dbName string) error {
	return migrateDataModelDb(ctx, db, dbName, getLatestMigration(ctx, db, dbName))
}

func RecreateDb(ctx context.Context, db *sql.DB, dbName string) error {
	if dbName != "test" {
		return configuration.LogAndReturnError(dbName, "", "Only test database can be recreated.")
	}
	rowsCount := 0
	if err := db.QueryRow("SELECT COUNT(uuid) FROM migrations").Scan(&rowsCount); err == nil && rowsCount >= 0 {
		configuration.Log(dbName, "", "Found %d rows in table.", rowsCount)
		commands := getDeletionSteps()
		keys := make([]int, 0)
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		var anyError error
		for _, step := range keys {
			if _, err := db.Exec(commands[step]); err != nil {
				configuration.LogError(dbName, "", "%d %q: %v", step, commands[step], err)
				anyError = err
			}
			configuration.Log(dbName, "", "Deletion %d executed.", step)
		}
		if anyError != nil {
			return anyError
		}
	}
	return migrateDb(ctx, db, dbName)
}

func getLatestMigration(ctx context.Context, db *sql.DB, dbName string) int {
	latestMigration := 0
	if err := db.QueryRow("SELECT MAX(int1) FROM migrations WHERE gridUuid = $1", model.UuidMigrations).Scan(&latestMigration); err != nil {
		configuration.Log(dbName, "", "No latest migration found: %v.", err)
		return 0
	}
	configuration.Log(dbName, "", "Latest migration: %d.", latestMigration)
	return latestMigration
}

func migrateDataModelDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int) error {
	commands := getMigrationSteps(dbName)
	keys := make([]int, 0)
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, step := range keys {
		err := migrateDbCommand(ctx, db, latestMigration, step, commands[step], dbName)
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateDbCommand(ctx context.Context, db *sql.DB, latestMigration int, step int, command string, dbName string) error {
	if step > latestMigration {
		_, err := db.Exec(command)
		if err != nil {
			return configuration.LogAndReturnError(dbName, "", "%d %q: %v", step, command, err)
		} else {
			newUuid := utils.GetNewUUID()
			_, err = db.Exec(getMigrationInsertStatement(), newUuid, step, command)
			if err != nil {
				return configuration.LogAndReturnError(dbName, "", "Can't insert into migrations: %v", err)
			}
			configuration.Log(dbName, "", "Migration %d executed.", step)
		}
	}
	return nil
}

// function is available for mocking
var getMigrationInsertStatement = func() string {
	return "INSERT INTO migrations " +
		"(uuid, " +
		"revision, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"int1, " +
		"text1) " +
		"VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"'" + model.UuidRootUser + "', " +
		"'" + model.UuidRootUser + "', " +
		"true, " +
		"'" + model.UuidMigrations + "', " +
		"$2, " +
		"$3)"
}

// function is available for mocking
var getRowsColumnDefinitions = func(grid *model.Grid) string {
	columnDefinitions := ""
	switch grid.Uuid {
	case model.UuidGrids:
		columnDefinitions += "text1, text2, text3, "
	case model.UuidColumns:
		columnDefinitions += "text1, text2, text3, int1, "
	case model.UuidRelationships:
		columnDefinitions += "text1, text2, text3, text4, text5, "
	default:
		for i := 1; i <= model.NumberOfTextFields; i++ {
			columnDefinitions += fmt.Sprintf("text%d text, ", i)
		}
		for i := 1; i <= model.NumberOfIntFields; i++ {
			columnDefinitions += fmt.Sprintf("int%d integer, ", i)
		}
	}
	return columnDefinitions
}
