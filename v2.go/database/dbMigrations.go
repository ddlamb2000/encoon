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

const (
	numberOfTextFields = 10
	numberOfIntFields  = 10
)

func migrateDb(ctx context.Context, db *sql.DB, dbName string) error {
	latestMigration, _ := migrateInitializationDb(ctx, db, dbName)
	return migrateDataModelDb(ctx, db, dbName, latestMigration)
}

func RecreateDb(ctx context.Context, db *sql.DB, dbName string) error {
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
	if err := db.QueryRow("SELECT MAX(int1) FROM rows WHERE gridUuid = $1", model.UuidMigrations).Scan(&latestMigration); err != nil {
		utils.Log("[%s] No latest migration found: %v.", dbName, err)
		return 0, nil
	}
	utils.Log("[%s] Latest migration: %d.", dbName, latestMigration)
	return latestMigration, nil
}

func migrateDataModelDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int) error {
	root, password := configuration.GetRootAndPassword(dbName)
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
				"VALUES ('" + model.UuidNumberColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumnTypes + "', " +
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
				"VALUES ('" + model.UuidUserColumnIdRel + "', " +
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
				"VALUES ('" + model.UuidUserColumnFirstNameRel + "', " +
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
				"VALUES ('" + model.UuidUserColumnLastNameRel + "', " +
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
				"VALUES ('" + model.UuidUserColumnPasswordRel + "', " +
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
				"VALUES ('" + model.UuidGridColumnUriRel + "', " +
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
				"VALUES ('" + model.UuidGridColumnNameRel + "', " +
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
				"VALUES ('" + model.UuidGridColumnDescRel + "', " +
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
				"VALUES ('" + model.UuidGridColumnIconRel + "', " +
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
				"VALUES ('" + model.UuidColumnColumnLabelRel + "', " +
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
				"VALUES ('" + model.UuidColumnColumnNameRel + "', " +
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
				"VALUES ('" + model.UuidColumnTypeColumnNameRel + "', " +
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
				"VALUES ('" + model.UuidColumnTypeColumnDescRel + "', " +
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
				"VALUES ('" + model.UuidMigrationsColumnSeqRel + "', " +
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
				"VALUES ('" + model.UuidMigrationsColumnCommandRel + "', " +
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
				"text2, " +
				"text3, " +
				"int1) " +
				"VALUES ('" + model.UuidColumnColumnColumnType + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Type', " +
				"'relationship1', " +
				"'_columntypes'," +
				"1)",

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
				"text2, " +
				"text3, " +
				"int1) " +
				"VALUES ('" + model.UuidGridColumnColumns + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Columns', " +
				"'relationship1', " +
				"'_columns'," +
				"99)",

			48: "INSERT INTO rows " +
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
				"VALUES ('" + model.UuidColumnColumnGridUri + "', " +
				"1, " +
				"NOW(), " +
				"NOW(), " +
				"'" + model.UuidRootUser + "', " +
				"'" + model.UuidRootUser + "', " +
				"true, " +
				"'" + model.UuidColumns + "', " +
				"'Grid uri', " +
				"'text3')",

			49: "INSERT INTO rows " +
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
				"VALUES ('" + model.UuidColumnColumnGridUriRel + "', " +
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
				"'" + model.UuidColumnColumnGridUri + "')",
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
				return utils.LogAndReturnError("[%s] Can't insert into migrations: %v", dbName, err)
			} else {
				utils.Log("[%s] Migration %v executed.", dbName, migration)
			}
		}
	}
	return nil
}

func getRowsColumnDefinitions() string {
	var columnDefinitions = ""
	for i := 1; i <= numberOfTextFields; i++ {
		columnDefinitions += fmt.Sprintf("text%d text, ", i)
	}
	for i := 1; i <= numberOfIntFields; i++ {
		columnDefinitions += fmt.Sprintf("int%d integer, ", i)
	}
	return columnDefinitions
}
