// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

// function is available for mocking
var getMigrationSteps = func(dbName string) map[int]string {
	root, password := configuration.GetRootAndPassword(dbName)
	return map[int]string{

		1: "CREATE TABLE rows (" +
			"gridUuid text NOT NULL, " +
			"uuid text NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			getRowsColumnDefinitions() +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid)" +
			")",

		2:  "CREATE INDEX rowsText1 ON rows(text1);",
		3:  "CREATE INDEX rowsText2 ON rows(text2);",
		4:  "CREATE INDEX rowsText3 ON rows(text3);",
		5:  "CREATE INDEX rowsText4 ON rows(text4);",
		6:  "CREATE INDEX rowsText5 ON rows(text5);",
		7:  "CREATE INDEX rowsText6 ON rows(text6);",
		8:  "CREATE INDEX rowsText7 ON rows(text7);",
		9:  "CREATE INDEX rowsText8 ON rows(text8);",
		10: "CREATE INDEX rowsText9 ON rows(text9);",
		11: "CREATE INDEX rowsText10 ON rows(text10);",

		100: "CREATE EXTENSION pgcrypto",

		107: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidGrids + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Grids', " +
			"'Data organized in rows and columns.', " +
			"'grid-3x3')",

		108: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidUsers + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Users', " +
			"'Users who has access to the system.', " +
			"'person')",

		109: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		110: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidColumns + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Columns', " +
			"'Columns of data grids.', " +
			"'columns')",

		111: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidMigrations + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Migrations', " +
			"'Statements run to create database.', " +
			"'file-earmark-text')",

		112: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidTransactions + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Transactions', " +
			"'Log of data changes.', " +
			"'journal')",

		113: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidColumnTypes + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Column types', " +
			"'Types of data grids columns.', " +
			"'columns-gap')",

		114: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		115: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		116: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		117: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidRelationships + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Relationships', " +
			"'Relationships between grid rows.', " +
			"'file-code')",

		118: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		119: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		120: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		121: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		122: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		123: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		124: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		125: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		127: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'text1')",

		128: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'text2')",

		129: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'text3')",

		131: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		132: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		133: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		134: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		135: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		136: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		137: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		138: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		139: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		140: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		141: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		142: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		143: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		144: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		145: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		146: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		147: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		150: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		151: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		152: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		153: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		154: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		156: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		157: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		158: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		159: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		160: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		162: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		163: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		164: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		165: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		166: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		167: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		168: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		169: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		170: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		171: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		172: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		173: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		174: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		175: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		176: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		177: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		178: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		179: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		180: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		181: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		182: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		183: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		184: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		185: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		186: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		187: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		188: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		189: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		190: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		191: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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

		199: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2, " +
			"text3) " +
			"VALUES ('" + model.UuidAccessLevel + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidGrids + "', " +
			"'Access level', " +
			"'Defines security access to data.', " +
			"'file-lock')",

		200: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2) " +
			"VALUES ('" + model.UuidAccessLevelColumnLevel + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidColumns + "', " +
			"'Level', " +
			"'text1')",

		201: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidAccessLevel + "', " +
			"'" + model.UuidColumns + "', " +
			"'" + model.UuidAccessLevelColumnLevel + "')",

		202: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidAccessLevelColumnLevel + "', " +
			"'" + model.UuidColumnTypes + "', " +
			"'" + model.UuidTextColumnType + "')",

		203: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1) " +
			"VALUES ('" + model.UuidAccessLevelReadAccess + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidAccessLevel + "', " +
			"'View access')",

		204: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1) " +
			"VALUES ('" + model.UuidAccessLevelWriteAccess + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidAccessLevel + "', " +
			"'Edit access')",

		205: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2) " +
			"VALUES ('" + model.UuidGridColumnAccessLevel + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidColumns + "', " +
			"'Default access', " +
			"'relationship2')",

		206: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnAccessLevel + "', " +
			"'" + model.UuidColumnTypes + "', " +
			"'" + model.UuidReferenceColumnType + "')",

		207: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnAccessLevel + "')",

		208: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnAccessLevel + "', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidAccessLevel + "')",

		209: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2) " +
			"VALUES ('" + model.UuidGridColumnOwner + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidColumns + "', " +
			"'Owner', " +
			"'relationship3')",

		210: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnOwner + "', " +
			"'" + model.UuidColumnTypes + "', " +
			"'" + model.UuidReferenceColumnType + "')",

		211: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnOwner + "')",

		212: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnOwner + "', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "')",

		213: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2) " +
			"VALUES ('" + model.UuidGridColumnReadAccessUsers + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidColumns + "', " +
			"'View access', " +
			"'relationship4')",

		214: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnReadAccessUsers + "', " +
			"'" + model.UuidColumnTypes + "', " +
			"'" + model.UuidReferenceColumnType + "')",

		215: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnReadAccessUsers + "')",

		216: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnReadAccessUsers + "', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "')",

		217: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1, " +
			"text2) " +
			"VALUES ('" + model.UuidGridColumnWriteAccessUsers + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidColumns + "', " +
			"'Edit access', " +
			"'relationship5')",

		218: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnWriteAccessUsers + "', " +
			"'" + model.UuidColumnTypes + "', " +
			"'" + model.UuidReferenceColumnType + "')",

		219: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnWriteAccessUsers + "')",

		220: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGridColumnWriteAccessUsers + "', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "')",

		221: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		222: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		223: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidColumns + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		224: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidMigrations + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		225: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidTransactions + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		226: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidColumnTypes + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		227: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidRelationships + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		228: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'relationship3', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidAccessLevel + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		229: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidAccessLevel + "', " +
			"'" + model.UuidAccessLevelSpecialAccess + "')",

		230: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
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
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidColumns + "', " +
			"'" + model.UuidAccessLevel + "', " +
			"'" + model.UuidAccessLevelSpecialAccess + "')",

		232: "INSERT INTO rows " +
			"(uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, " +
			"text1) " +
			"VALUES ('" + model.UuidAccessLevelSpecialAccess + "', " +
			"1, " +
			"NOW(), " +
			"NOW(), " +
			"'" + model.UuidRootUser + "', " +
			"'" + model.UuidRootUser + "', " +
			"true, " +
			"'" + model.UuidAccessLevel + "', " +
			"'Special access')",
	}
}

// function is available for mocking
var getDeletionSteps = func() map[int]string {
	return map[int]string{
		1: "DROP TABLE rows",
		2: "DROP EXTENSION pgcrypto",
	}
}
