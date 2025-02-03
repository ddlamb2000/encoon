// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

// function is available for mocking
var getMigrationSteps = func(dbName string) map[int]string {
	root, password := configuration.GetRootAndPassword(dbName)
	return map[int]string{

		1: "CREATE TABLE migrations (" +
			"gridUuid uuid NOT NULL, " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
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
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
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
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
			"enabled boolean NOT NULL" +
			getRowsColumnDefinitions(model.GetNewGrid("")) + ", " +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid)" +
			")",

		30: "CREATE TABLE columns (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"text2 text," +
			"text3 text," +
			"int1 integer," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",

		40: "CREATE TABLE relationships (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
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
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"text2 text," +
			"text3 text," +
			"text4 text," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",

		55: "CREATE TABLE transactions (" +
			"gridUuid uuid NOT NULL REFERENCES grids (uuid), " +
			"uuid uuid NOT NULL, " +
			"created timestamp with time zone NOT NULL, " +
			"createdBy uuid NOT NULL, " +
			"updated timestamp with time zone NOT NULL, " +
			"updatedBy uuid NOT NULL, " +
			"enabled boolean NOT NULL, " +
			"text1 text," +
			"revision integer NOT NULL CHECK (revision > 0), " +
			"PRIMARY KEY (gridUuid, uuid), " +
			"UNIQUE (uuid)" +
			")",

		100: "CREATE EXTENSION IF NOT EXISTS pgcrypto",

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

		110: "CREATE EXTENSION IF NOT EXISTS vector",

		111: "ALTER TABLE rows ADD COLUMN embedding vector(768)",

		112: "ALTER TABLE rows ADD COLUMN revisionEmbedding integer NOT NULL DEFAULT 0",

		113: "ALTER TABLE rows ADD COLUMN tokenCount integer NOT NULL DEFAULT 0",

		114: "ALTER TABLE rows ADD COLUMN embeddingText text NOT NULL DEFAULT ''",
	}
}

// function is available for mocking
var getDeletionSteps = func() map[int]string {
	return map[int]string{
		1: "DROP EXTENSION pgcrypto",
		2: "DROP TABLE transactions",
		3: "DROP TABLE users",
		4: "DROP TABLE relationships",
		5: "DROP TABLE columns",
		6: "DROP TABLE rows",
		7: "DROP TABLE grids",
		8: "DROP TABLE migrations",
	}
}
