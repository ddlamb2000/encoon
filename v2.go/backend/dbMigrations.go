// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"database/sql"
	"sort"

	"d.lambert.fr/encoon/utils"
)

func migrateDb(ctx context.Context, db *sql.DB, dbName string) error {
	latestMigration, _ := migrateInitializationDb(ctx, db, dbName)
	err := migrateDataModelDb(ctx, db, dbName, latestMigration)
	if err != nil {
		return err
	}
	return nil
}

func recreateDb(ctx context.Context, db *sql.DB, dbName string) error {
	if dbName != "test" {
		return utils.LogAndReturnError("[%s] Only test database can be recreated.", dbName)
	}
	var rowsCount int
	if err := db.QueryRow("SELECT COUNT(uuid) FROM rows").Scan(&rowsCount); err == nil && rowsCount >= 0 {
		utils.Log("[%s] Found %d rows in table.", dbName, rowsCount)
		command := "DROP TABLE rows"
		if _, err := db.Exec(command); err != nil {
			return utils.LogAndReturnError("[%s] Can't delete database rows: %v", dbName, err)
		}
		command = "DROP EXTENSION pgcrypto"
		if _, err := db.Exec(command); err != nil {
			return utils.LogAndReturnError("[%s] Can't drop extension: %v", dbName, err)
		}
		utils.Log("[%s] Table rows removed.", dbName)
	}
	return migrateDb(ctx, db, dbName)
}

func migrateInitializationDb(ctx context.Context, db *sql.DB, dbName string) (int, error) {
	var latestMigration int = 0
	if err := db.QueryRow("SELECT MAX(int01) FROM rows WHERE gridUuid = $1", utils.UuidMigrations).Scan(&latestMigration); err != nil {
		utils.Log("[%s] No latest migration found: %v.", dbName, err)
		return 0, nil
	}
	utils.Log("[%s] Latest migration: %d.", dbName, latestMigration)
	return latestMigration, nil
}

func migrateDataModelDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int) error {
	root, password := utils.GetRootAndPassword(dbName)
	err := migrateCommandsDb(ctx, db, dbName, latestMigration,
		map[int]string{
			1: "CREATE TABLE rows (" +
				"uuid uuid NOT NULL PRIMARY KEY, " +
				"created timestamp with time zone, " +
				"createdBy uuid, " +
				"updated timestamp with time zone, " +
				"updatedBy uuid, " +
				"enabled boolean, " +
				"gridUuid uuid, " +
				"parentUuid uuid, " +
				getRowsColumnDefinitions() +
				"version integer)",

			2: "CREATE EXTENSION pgcrypto",

			4: "CREATE INDEX gridUuid ON rows(gridUuid);",

			5: "CREATE INDEX gridParentUuid ON rows(parentUuid);",

			6: "CREATE INDEX gridText01 ON rows(text01);",

			7: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidGrids + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidGrids + "', " +
				"'_grids', " +
				"'Grids', " +
				"'Data organized in rows and columns.', " +
				"'grid-3x3')",

			8: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidUsers + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidGrids + "', " +
				"'_users', " +
				"'Users', " +
				"'Users who has access to the system.', " +
				"'person')",

			9: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidRootUser + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidUsers + "', " +
				"'" + root + "', " +
				"'" + root + "', " +
				"'" + root + "', " +
				"'" + password + "')",

			10: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidColumns + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidGrids + "', " +
				"'_columns', " +
				"'Columns', " +
				"'Columns of data grids.', " +
				"'columns')",

			11: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidMigrations + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidGrids + "', " +
				"'_migrations', " +
				"'Migrations', " +
				"'Statements run to create database.', " +
				"'file-earmark-text')",

			12: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidTransactions + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidGrids + "', " +
				"'_transactions', " +
				"'Transactions', " +
				"'Log of data changes.', " +
				"'journal')",

			13: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02, " +
				"text03, " +
				"text04) " +
				"VALUES ('" + utils.UuidColumnTypes + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidGrids + "', " +
				"'_columntypes', " +
				"'Column types', " +
				"'Types of data grids columns.', " +
				"'columns-gap')",

			14: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02) " +
				"VALUES ('" + utils.UuidTextColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidColumnTypes + "', " +
				"'Text', " +
				"'Text column type.')",

			15: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02) " +
				"VALUES ('" + utils.UuidNumberColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidColumnTypes + "', " +
				"'Number', " +
				"'Number column type.')",

			16: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " +
				"text02) " +
				"VALUES ('" + utils.UuidReferenceColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidColumnTypes + "', " +
				"'Reference', " +
				"'Reference to another data grid row column type.')",
		})
	if err != nil {
		return err
	}
	return nil
}

func migrateCommandsDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int, commands map[int]string) error {
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

func migrateDbCommand(ctx context.Context, db *sql.DB, latestMigration int, migration int, command string, dbName string) error {
	if migration > latestMigration {
		_, err := db.Exec(command)
		if err != nil {
			return utils.LogAndReturnError("[%s] %d %q: %v", dbName, migration, command, err)
		} else {
			insertMigrationStatement := "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"int01, " +
				"text01) " +
				"VALUES ($1, " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + utils.UuidRootUser + "', " +
				"'" + utils.UuidRootUser + "', " +
				"true, " +
				"'" + utils.UuidMigrations + "', " +
				"$2, " +
				"$3)"
			newUuid := utils.GetNewUUID()
			_, err = db.Exec(insertMigrationStatement, newUuid, migration, command)
			if err != nil {
				return utils.LogAndReturnError("[%s] Can't insert into migrations: %v", dbName, err)
			} else {
				utils.Log("[%s] Migration %v executed.", dbName, migration)
			}
		}
	}
	return nil
}
