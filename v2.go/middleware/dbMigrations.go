// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/backend/utils"
)

func migrateDb(ctx context.Context, db *sql.DB, dbName string) {
	var latestMigration int
	if err := db.QueryRow("select 1 from migrations").Scan(&latestMigration); err != nil {
		if err == sql.ErrNoRows {
			utils.Log("[%q] Migration table exists.", dbName)
		} else {
			utils.Log("[%q] Migrations table doesn't exist: %v.", dbName, err)
			command := "CREATE TABLE migrations (migration integer, command text)"
			_, err := db.Exec(command)
			if err != nil {
				utils.LogError("[%q] Create table migrations: %v", dbName, err)
				return
			}
			utils.Log("[%q] Migration table created.", dbName)
			_, err = db.Exec("INSERT INTO migrations (migration, command) VALUES ($1, $2)", 1, command)
			if err != nil {
				utils.LogError("[%q] Insert into migrations: %v", dbName, err)
			}
		}
	}

	if err := db.QueryRow("select max(migration) from migrations").Scan(&latestMigration); err != nil {
		utils.LogError("[%q] Can't access migrations table %v.", dbName, err)
		return
	} else {
		utils.Log("[%q] Latest migration: %v.", dbName, latestMigration)
		migrageDbCommand(ctx, db, latestMigration, 2, "CREATE TABLE users (uuid text, version integer, enabled boolean, id text, firstName text, lastName text, password text)", dbName)
		root, password := utils.GetRootAndPassowrd(dbName)
		migrageDbCommand(ctx, db, latestMigration, 3, "INSERT INTO users (uuid, version, enabled, id, password) VALUES ('0', 1, true, '"+root+"', '"+password+"')", dbName)
	}
}

func migrageDbCommand(ctx context.Context, db *sql.DB, latestMigration int, migration int, command string, dbName string) {
	if migration > latestMigration {
		_, err := db.Exec(command)
		if err != nil {
			utils.LogError("[%q] Command %q command: %v", dbName, command, err)
		} else {
			_, err = db.Exec("INSERT INTO migrations (migration, command) VALUES ($1, $2)", migration, command)
			if err != nil {
				utils.LogError("[%q] Insert into migrations: %v", dbName, err)
			} else {
				utils.Log("[%q] Latest migration: %v.", dbName, migration)
			}
		}
	}
}
