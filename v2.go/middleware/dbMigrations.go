// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"
	"sort"

	"d.lambert.fr/encoon/utils"
)

func migrateDb(ctx context.Context, db *sql.DB, dbName string) {
	if latestMigration, error := migrateInitializationDb(ctx, db, dbName); !error && latestMigration > 0 {
		migrateDataModelDb(ctx, db, dbName, latestMigration)
	}
}

func migrateInitializationDb(ctx context.Context, db *sql.DB, dbName string) (int, bool) {
	var latestMigration int
	if err := db.QueryRow("SELECT 1 FROM migrations").Scan(&latestMigration); err != nil {
		if err == sql.ErrNoRows {
			utils.Log("[%q] Migration table exists.", dbName)
		} else {
			utils.Log("[%q] Migrations table doesn't exist: %v.", dbName, err)
			command := "CREATE TABLE migrations (migration integer, command text)"
			_, err := db.Exec(command)
			if err != nil {
				utils.LogError("[%q] Create table migrations: %v", dbName, err)
				return latestMigration, true
			}
			utils.Log("[%q] Migration table created.", dbName)
			_, err = db.Exec("INSERT INTO migrations (migration, command) VALUES ($1, $2)", 1, command)
			if err != nil {
				utils.LogError("[%q] Insert into migrations: %v", dbName, err)
				return latestMigration, true
			}
		}
	}

	if err := db.QueryRow("SELECT MAX(migration) FROM migrations").Scan(&latestMigration); err != nil {
		utils.LogError("[%q] Can't access migrations table %v.", dbName, err)
		return latestMigration, true
	}
	return latestMigration, false
}

func migrateDataModelDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int) {
	root, password := utils.GetRootAndPassword(dbName)

	migrateCommandsDb(ctx, db, dbName, latestMigration,
		map[int]string{
			2: "CREATE EXTENSION pgcrypto",

			3: "CREATE TABLE rows (" +
				"uuid uuid NOT NULL PRIMARY KEY, " +
				"version integer, " +
				"created timestamp with time zone, " +
				"createdBy uuid, " +
				"updated timestamp with time zone, " +
				"updatedBy uuid, " +
				"enabled boolean, " +
				"gridUuid uuid, " +
				"parentUuid uuid, " +
				"text01 text, " +
				"text02 text, " +
				"text03 text, " +
				"text04 text, " +
				"text05 text, " +
				"text06 text, " +
				"text07 text, " +
				"text08 text, " +
				"text09 text, " +
				"text10 text)",

			4: "CREATE INDEX gridUuid ON rows(gridUuid);",

			5: "CREATE INDEX gridParentUuid ON rows(parentUuid);",

			6: "CREATE INDEX gridText01 ON rows(text01);",

			7: "CREATE INDEX gridText02 ON rows(text02);",

			8: "CREATE INDEX gridText03 ON rows(text03);",

			9: "CREATE INDEX gridText04 ON rows(text04);",

			10: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01) " +
				"VALUES ('f35ef7de-66e7-4e51-9a09-6ff8667da8f7', " + // Grid: Grids
				"1, " +
				"'January 8 04:05:06 1999 PST', " +
				"'January 8 04:05:06 1999 PST', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"true, " +
				"'f35ef7de-66e7-4e51-9a09-6ff8667da8f7', " +
				"'grids')",

			11: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01) " +
				"VALUES ('018803e1-b4bf-42fa-b58f-ac5faaeeb0c2', " + // Grid: Users
				"1, " +
				"'January 8 04:05:06 1999 PST', " +
				"'January 8 04:05:06 1999 PST', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"true, " +
				"'f35ef7de-66e7-4e51-9a09-6ff8667da8f7', " +
				"'users')",

			12: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01, " + // id
				"text02, " + // firstName
				"text03, " + // lastName
				"text04) " + // password
				"VALUES ('3a33485c-7683-4482-aa5d-0aa51e58d79d', " + // Users: root
				"1, " +
				"'January 8 04:05:06 1999 PST', " +
				"'January 8 04:05:06 1999 PST', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"true, " +
				"'018803e1-b4bf-42fa-b58f-ac5faaeeb0c2', " +
				"'" + root + "', " +
				"'" + root + "', " +
				"'" + root + "', " +
				"'" + password + "')",

			13: "INSERT INTO rows " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"createdBy, " +
				"updatedBy, " +
				"enabled, " +
				"gridUuid, " +
				"text01) " +
				"VALUES ('533b6862-add3-4fef-8f93-20a17aaaaf5a', " + // Grid: Columns
				"1, " +
				"'January 8 04:05:06 1999 PST', " +
				"'January 8 04:05:06 1999 PST', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"'3a33485c-7683-4482-aa5d-0aa51e58d79d', " +
				"true, " +
				"'f35ef7de-66e7-4e51-9a09-6ff8667da8f7', " +
				"'columns')",
		})
}

func migrateCommandsDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int, commands map[int]string) {
	keys := make([]int, 0)
	i := 0
	for k, _ := range commands {
		keys = append(keys, k)
		i++
	}
	sort.Ints(keys)
	for _, step := range keys {
		migrateDbCommand(ctx, db, latestMigration, step, commands[step], dbName)
	}
}

func migrateDbCommand(ctx context.Context, db *sql.DB, latestMigration int, migration int, command string, dbName string) {
	if migration > latestMigration {
		_, err := db.Exec(command)
		if err != nil {
			utils.LogError("[%q] %d %q: %v", dbName, migration, command, err)
		} else {
			_, err = db.Exec("INSERT INTO migrations (migration, command) VALUES ($1, $2)", migration, command)
			if err != nil {
				utils.LogError("[%q] Insert into migrations: %v", dbName, err)
			} else {
				utils.Log("[%q] Migration %v executed.", dbName, migration)
			}
		}
	}
}
