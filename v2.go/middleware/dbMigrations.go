// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/backend/utils"
)

func migrateDb(ctx context.Context, db *sql.DB) {
	var count int
	if err := db.QueryRow("select 1 from migrations").Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			utils.Log("Migration table exists.")
		} else {
			utils.Log("Migrations table doesn't exist: %v.", err)
			_, err := db.Exec("CREATE TABLE migrations (migration integer, command text)")
			if err != nil {
				utils.LogError("Create table migrations: %v", err)
				return
			}
			utils.Log("Migration table created.")
			_, err = db.Exec("INSERT INTO migrations (migration, command) VALUES (1, 'Migrations table initialized')")
			if err != nil {
				utils.LogError("Insert into migrations: %v", err)
				return
			}
		}
	}

	if err := db.QueryRow("select max(migration) from migrations").Scan(&count); err != nil {
		utils.LogError("Can't access migrations table %v.", err)
		return
	} else {
		utils.Log("Latest migration: %v.", count)
	}
}
