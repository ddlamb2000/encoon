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
	latestMigration, _ := migrateInitializationDb(ctx, db, dbName)
	return migrateDataModelDb(ctx, db, dbName, latestMigration)
}

func RecreateDb(ctx context.Context, db *sql.DB, dbName string) error {
	if dbName != "test" {
		return configuration.LogAndReturnError("[%s] Only test database can be recreated.", dbName)
	}
	var rowsCount int
	if err := db.QueryRow("SELECT COUNT(uuid) FROM rows").Scan(&rowsCount); err == nil && rowsCount >= 0 {
		configuration.Log("[%s] Found %d rows in table.", dbName, rowsCount)
		command := "DROP TABLE rows"
		if _, err := db.Exec(command); err != nil {
			return configuration.LogAndReturnError("[%s] Can't delete database rows: %v", dbName, err)
		}
		command = "DROP EXTENSION pgcrypto"
		if _, err := db.Exec(command); err != nil {
			return configuration.LogAndReturnError("[%s] Can't drop extension: %v", dbName, err)
		}
		configuration.Log("[%s] Table rows removed.", dbName)
	}
	return migrateDb(ctx, db, dbName)
}

func migrateInitializationDb(ctx context.Context, db *sql.DB, dbName string) (int, error) {
	var latestMigration int = 0
	if err := db.QueryRow("SELECT MAX(int1) FROM rows WHERE gridUuid = $1", model.UuidMigrations).Scan(&latestMigration); err != nil {
		configuration.Log("[%s] No latest migration found: %v.", dbName, err)
		return 0, nil
	}
	configuration.Log("[%s] Latest migration: %d.", dbName, latestMigration)
	return latestMigration, nil
}

func migrateDataModelDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int) error {
	root, password := configuration.GetRootAndPassword(dbName)
	err := migrateCommandsDb(ctx, db, dbName, latestMigration,
		map[int]string{
			1: "CREATE TABLE rows (" +
				"uuid text NOT NULL PRIMARY KEY, " +
				"created timestamp with time zone, " +
				"createdBy text, " +
				"updated timestamp with time zone, " +
				"updatedBy text, " +
				"enabled boolean, " +
				"gridUuid text, " +
				getRowsColumnDefinitions() +
				"version integer)",

			2: "CREATE EXTENSION pgcrypto",
			4: "CREATE INDEX gridUuid ON rows(gridUuid);",
			5: "CREATE INDEX gridText1 ON rows(text1);",
			6: "CREATE INDEX gridText2 ON rows(text2);",

			7: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidGrids + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
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
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidUsers + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
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
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidRootUser + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidUsers + "', " +
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
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidColumns + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
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
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidMigrations + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
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
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidTransactions + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
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
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidColumnTypes + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
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
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidTextColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumnTypes + "', " +
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
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidIntColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumnTypes + "', " +
				"'Integer', " +
				"'Integer number column type.')",

			16: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidReferenceColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumnTypes + "', " +
				"'Reference', " +
				"'Reference to another data grid row column type.')",

			17: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4) " +
				"VALUES ('" + model.UuidRelationships + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidGrids + "', " +
				"'_relationships', " +
				"'Relationships', " +
				"'Relationships between grid rows.', " +
				"'file-code')",

			18: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidUserColumnId + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Username', " +
				"'text1')",

			19: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidUserColumnFirstName + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'First name', " +
				"'text2')",

			20: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidUserColumnLastName + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Last name', " +
				"'text3')",

			21: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidUserColumnPassword + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Password', " +
				"'text4')",

			22: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidUsers + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnId + "')",

			23: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidUsers + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnFirstName + "')",

			24: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidUsers + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnLastName + "')",

			25: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidUsers + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnPassword + "')",

			26: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidGridColumnUri + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Uri', " +
				"'text1')",

			27: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidGridColumnName + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Name', " +
				"'text2')",

			28: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidGridColumnDesc + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Description', " +
				"'text3')",

			29: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidGridColumnIcon + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Icon', " +
				"'text4')",

			30: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnUri + "')",

			31: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnName + "')",

			32: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnDesc + "')",

			33: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnIcon + "')",

			34: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidColumnColumnLabel + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Label', " +
				"'text1')",

			35: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidColumnColumnName + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Name', " +
				"'text2')",

			36: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnLabel + "')",

			37: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnName + "')",

			38: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidColumnTypeColumnName + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Name', " +
				"'text1')",

			39: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidColumnTypeColumnDesc + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Description', " +
				"'text2')",

			40: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnName + "')",

			41: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnDesc + "')",

			42: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidMigrationsColumnSeq + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Sequence', " +
				"'int1')",

			43: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidMigrationsColumnCommand + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Command', " +
				"'text1')",

			44: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidMigrations + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidMigrationsColumnSeq + "')",

			45: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidMigrations + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidMigrationsColumnCommand + "')",

			46: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidColumnColumnColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Type', " +
				"'relationship1')",

			47: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidGridColumnColumns + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Columns', " +
				"'relationship1')",

			50: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnId + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			51: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnFirstName + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			52: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnLastName + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			53: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidPasswordColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumnTypes + "', " +
				"'Password', " +
				"'Password or passphrase.')",

			54: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidUserColumnPassword + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidPasswordColumnType + "')",

			55: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnUri + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			56: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnName + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			57: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnDesc + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			58: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnIcon + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			59: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnLabel + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			60: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnName + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			62: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnColumnType + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidReferenceColumnType + "')",

			63: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnName + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			64: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnDesc + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			65: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidMigrationsColumnSeq + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidIntColumnType + "')",

			66: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidMigrationsColumnCommand + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			67: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnColumnType + "')",

			68: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnColumns + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidReferenceColumnType + "')",

			69: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnColumns + "')",

			70: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidUuidColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumnTypes + "', " +
				"'Uuid', " +
				"'Universally unique identifier.')",

			71: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidRelationshipColumnName + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Name', " +
				"'text1')",

			72: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidRelationships + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnName + "')",

			73: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnName + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidTextColumnType + "')",

			74: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidRelationshipColumnFromGridUuid + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'From grid', " +
				"'text2')",

			75: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidRelationships + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnFromGridUuid + "')",

			76: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnFromGridUuid + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidUuidColumnType + "')",

			77: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidRelationshipColumnFromUuid + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'From row', " +
				"'text3')",

			78: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidRelationships + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnFromUuid + "')",

			79: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnFromUuid + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidUuidColumnType + "')",

			80: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidRelationshipColumnToGridUuid + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'To grid', " +
				"'text4')",

			81: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidRelationships + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnToGridUuid + "')",

			82: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnToGridUuid + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidUuidColumnType + "')",

			83: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidRelationshipColumnToUuid + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'To row', " +
				"'text5')",

			84: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidRelationships + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnToUuid + "')",

			85: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidRelationshipColumnToUuid + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidUuidColumnType + "')",

			86: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2) " +
				"VALUES ('" + model.UuidColumnTypeColumnGridPrompt + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Grid for prompt', " +
				"'relationship2')",

			87: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnGridPrompt + "')",

			88: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship1', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnGridPrompt + "', " +
				"'" + model.UuidColumnTypes + "', " +
				"'" + model.UuidReferenceColumnType + "')",

			89: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship2', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnColumnColumnType + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumnTypes + "')",

			90: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship2', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidGridColumnColumns + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidColumns + "')",

			91: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text1, " +
				"text2, " +
				"text3, " +
				"text4, " +
				"text5) " +
				"VALUES ('" + utils.GetNewUUID() + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidRelationships + "', " +
				"'relationship2', " +
				"'" + model.UuidColumns + "', " +
				"'" + model.UuidColumnTypeColumnGridPrompt + "', " +
				"'" + model.UuidGrids + "', " +
				"'" + model.UuidGrids + "')",
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
			return configuration.LogAndReturnError("[%s] %d %q: %v", dbName, migration, command, err)
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
			newUuid := utils.GetNewUUID()
			_, err = db.Exec(insertMigrationStatement, newUuid, migration, command)
			if err != nil {
				return configuration.LogAndReturnError("[%s] Can't insert into migrations: %v", dbName, err)
			} else {
				configuration.Log("[%s] Migration %v executed.", dbName, migration)
			}
		}
	}
	return nil
}

func getRowsColumnDefinitions() string {
	var columnDefinitions = ""
	for i := 1; i <= model.NumberOfTextFields; i++ {
		columnDefinitions += fmt.Sprintf("text%d text, ", i)
	}
	for i := 1; i <= model.NumberOfIntFields; i++ {
		columnDefinitions += fmt.Sprintf("int%d integer, ", i)
	}
	return columnDefinitions
}
