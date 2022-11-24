// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

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
	var rowsCount int
	if err := db.QueryRow("SELECT COUNT(uuid) FROM rows").Scan(&rowsCount); err == nil && rowsCount >= 0 {
		configuration.Log(dbName, "", "Found %d rows in table.", rowsCount)
		commands := getDeletionSteps()
		keys := make([]int, 0)
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, step := range keys {
			if _, err := db.Exec(commands[step]); err != nil {
				return configuration.LogAndReturnError(dbName, "", "%d %q: %v", step, commands[step], err)
			}
			configuration.Log(dbName, "", "Deletion %d executed.", step)
		}
	}
	return migrateDb(ctx, db, dbName)
}

func getLatestMigration(ctx context.Context, db *sql.DB, dbName string) int {
	var latestMigration int = 0
	if err := db.QueryRow("SELECT MAX(int1) FROM rows WHERE gridUuid = $1", model.UuidMigrations).Scan(&latestMigration); err != nil {
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
		configuration.Log(dbName, "", "Migration %d executed.", step)
	}
	return nil
}

func migrateDbCommand(ctx context.Context, db *sql.DB, latestMigration int, migration int, command string, dbName string) error {
	if migration > latestMigration {
		_, err := db.Exec(command)
		if err != nil {
			return configuration.LogAndReturnError(dbName, "", "%d %q: %v", migration, command, err)
		} else {
			newUuid := utils.GetNewUUID()
			_, err = db.Exec(getMigrationInsertStatement(), newUuid, migration, command)
			if err != nil {
				return configuration.LogAndReturnError(dbName, "", "Can't insert into migrations: %v", err)
			}
		}
	}
	return nil
}

// function is available for mocking
var getMigrationInsertStatement = func() string {
	return "INSERT INTO rows " +
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
var getRowsColumnDefinitions = func() string {
	var columnDefinitions = ""
	for i := 1; i <= model.NumberOfTextFields; i++ {
		columnDefinitions += fmt.Sprintf("text%d text, ", i)
	}
	for i := 1; i <= model.NumberOfIntFields; i++ {
		columnDefinitions += fmt.Sprintf("int%d integer, ", i)
	}
	return columnDefinitions
}
