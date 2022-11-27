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

		1: "CREATE TABLE migrations (" +
			"gridUuid uuid NOT NULL, " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"int1 integer," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid)" +
			")",

		5: "CREATE TABLE grids (" +
			"gridUuid uuid NOT NULL, " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"text2 text," +
			"text3 text," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",

		10: "CREATE TABLE rows (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			getRowsColumnDefinitions() +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid)" +
			")",

		30: "CREATE TABLE columns (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"text2 text," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",

		40: "CREATE TABLE relationships (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"text2 uuid REFERENCES grids (uuid)," +
			"text3 uuid," +
			"text4 uuid REFERENCES grids (uuid)," +
			"text5 uuid," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",
		41: "CREATE INDEX relationships_text1 ON relationships (text1)",
		42: "CREATE INDEX relationships_text2 ON relationships (text2)",
		43: "CREATE INDEX relationships_text3 ON relationships (text3)",
		44: "CREATE INDEX relationships_text4 ON relationships (text4)",
		45: "CREATE INDEX relationships_text5 ON relationships (text5)",

		50: "CREATE TABLE users (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy text NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy text NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"text2 text," +
			"text3 text," +
			"text4 text," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",

		100: "CREATE EXTENSION pgcrypto",

		107: "INSERT INTO grids " +
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

		108: "INSERT INTO grids " +
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

		109: "INSERT INTO users " +
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

		110: "INSERT INTO grids " +
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

		111: "INSERT INTO grids " +
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

		112: "INSERT INTO grids " +
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

		113: "INSERT INTO grids " +
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

		117: "INSERT INTO grids " +
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

		118: "INSERT INTO columns " +
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

		119: "INSERT INTO columns " +
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

		120: "INSERT INTO columns " +
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

		121: "INSERT INTO columns " +
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

		122: "INSERT INTO relationships " +
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

		123: "INSERT INTO relationships " +
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

		124: "INSERT INTO relationships " +
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

		125: "INSERT INTO relationships " +
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

		127: "INSERT INTO columns " +
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

		128: "INSERT INTO columns " +
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

		129: "INSERT INTO columns " +
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

		131: "INSERT INTO relationships " +
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

		132: "INSERT INTO relationships " +
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

		133: "INSERT INTO relationships " +
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

		134: "INSERT INTO columns " +
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

		135: "INSERT INTO columns " +
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

		136: "INSERT INTO relationships " +
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

		137: "INSERT INTO relationships " +
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

		138: "INSERT INTO columns " +
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

		139: "INSERT INTO columns " +
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

		140: "INSERT INTO relationships " +
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

		141: "INSERT INTO relationships " +
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

		142: "INSERT INTO columns " +
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

		143: "INSERT INTO columns " +
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

		144: "INSERT INTO relationships " +
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

		145: "INSERT INTO relationships " +
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

		146: "INSERT INTO columns " +
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

		147: "INSERT INTO columns " +
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

		150: "INSERT INTO relationships " +
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

		151: "INSERT INTO relationships " +
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

		152: "INSERT INTO relationships " +
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

		154: "INSERT INTO relationships " +
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

		156: "INSERT INTO relationships " +
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

		157: "INSERT INTO relationships " +
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

		158: "INSERT INTO relationships " +
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

		159: "INSERT INTO relationships " +
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

		160: "INSERT INTO relationships " +
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

		162: "INSERT INTO relationships " +
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

		163: "INSERT INTO relationships " +
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

		164: "INSERT INTO relationships " +
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

		165: "INSERT INTO relationships " +
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

		166: "INSERT INTO relationships " +
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

		167: "INSERT INTO relationships " +
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

		168: "INSERT INTO relationships " +
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

		169: "INSERT INTO relationships " +
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

		171: "INSERT INTO columns " +
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

		172: "INSERT INTO relationships " +
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

		173: "INSERT INTO relationships " +
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

		174: "INSERT INTO columns " +
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

		175: "INSERT INTO relationships " +
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

		176: "INSERT INTO relationships " +
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

		177: "INSERT INTO columns " +
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

		178: "INSERT INTO relationships " +
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

		179: "INSERT INTO relationships " +
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

		180: "INSERT INTO columns " +
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

		181: "INSERT INTO relationships " +
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

		182: "INSERT INTO relationships " +
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

		183: "INSERT INTO columns " +
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

		184: "INSERT INTO relationships " +
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

		185: "INSERT INTO relationships " +
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

		186: "INSERT INTO columns " +
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

		187: "INSERT INTO relationships " +
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

		188: "INSERT INTO relationships " +
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

		189: "INSERT INTO relationships " +
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

		190: "INSERT INTO relationships " +
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

		191: "INSERT INTO relationships " +
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

		199: "INSERT INTO grids " +
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

		200: "INSERT INTO columns " +
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

		201: "INSERT INTO relationships " +
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

		202: "INSERT INTO relationships " +
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

		205: "INSERT INTO columns " +
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

		206: "INSERT INTO relationships " +
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

		207: "INSERT INTO relationships " +
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

		208: "INSERT INTO relationships " +
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

		209: "INSERT INTO columns " +
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
			"'Owners', " +
			"'relationship3')",

		210: "INSERT INTO relationships " +
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

		211: "INSERT INTO relationships " +
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

		212: "INSERT INTO relationships " +
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

		213: "INSERT INTO columns " +
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

		214: "INSERT INTO relationships " +
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

		215: "INSERT INTO relationships " +
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

		216: "INSERT INTO relationships " +
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

		217: "INSERT INTO columns " +
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

		218: "INSERT INTO relationships " +
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

		219: "INSERT INTO relationships " +
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

		220: "INSERT INTO relationships " +
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

		221: "INSERT INTO relationships " +
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

		222: "INSERT INTO relationships " +
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

		223: "INSERT INTO relationships " +
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

		224: "INSERT INTO relationships " +
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

		225: "INSERT INTO relationships " +
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

		226: "INSERT INTO relationships " +
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

		227: "INSERT INTO relationships " +
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

		228: "INSERT INTO relationships " +
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

		233: "INSERT INTO relationships " +
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
			"'relationship4', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",

		234: "INSERT INTO relationships " +
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
			"'relationship5', " +
			"'" + model.UuidGrids + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidUsers + "', " +
			"'" + model.UuidRootUser + "')",
	}
}

// function is available for mocking
var getDeletionSteps = func() map[int]string {
	return map[int]string{
		7: "DROP TABLE migrations",
		6: "DROP TABLE grids",
		5: "DROP TABLE rows",
		4: "DROP TABLE columns",
		3: "DROP TABLE relationships",
		2: "DROP TABLE users",
		1: "DROP EXTENSION pgcrypto",
	}
}
