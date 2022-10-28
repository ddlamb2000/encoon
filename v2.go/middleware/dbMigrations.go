// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/backend/utils"
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
	root, password := utils.GetRootAndPassowrd(dbName)

	migrateCommandsDb(ctx, db, dbName, latestMigration,
		map[int]string{
			2: "CREATE TABLE users (" +
				"uuid uuid, " +
				"version integer, " +
				"created timestamp with time zone, " +
				"updated timestamp with time zone, " +
				"enabled boolean, " +
				"id text, " +
				"firstName text, " +
				"lastName text, " +
				"password text)",

			3: "INSERT INTO users " +
				"(uuid, " +
				"version, " +
				"created, " +
				"updated, " +
				"enabled, " +
				"id, " +
				"firstName, " +
				"lastName, " +
				"password) " +
				"VALUES ('3a33485c-7683-4482-aa5d-0aa51e58d79d'," +
				" 1, " +
				"'January 8 04:05:06 1999 PST', " +
				"'January 8 04:05:06 1999 PST', " +
				"true, " +
				"'" + root + "', " +
				"'" + root + "', " +
				"'" + root + "', " +
				"'" + password + "')",
		})
}

func migrateCommandsDb(ctx context.Context, db *sql.DB, dbName string, latestMigration int, commands map[int]string) {
	for step, command := range commands {
		migrateDbCommand(ctx, db, latestMigration, step, command, dbName)
	}
}

func migrateDbCommand(ctx context.Context, db *sql.DB, latestMigration int, migration int, command string, dbName string) {
	if migration > latestMigration {
		_, err := db.Exec(command)
		if err != nil {
			utils.LogError("[%q] Command %q command: %v", dbName, command, err)
		} else {
			_, err = db.Exec("INSERT INTO migrations (migration, command) VALUES ($1, $2)", migration, command)
			if err != nil {
				utils.LogError("[%q] Insert into migrations: %v", dbName, err)
			} else {
				utils.Log("[%q] Migration: %v.", dbName, migration)
			}
		}
	}
}
