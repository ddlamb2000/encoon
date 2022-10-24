// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"d.lambert.fr/encoon/backend/utils"
)

var dbs = make(map[string]*sql.DB)

func ConnectDbServers(dbNames []string) {
	for _, dbName := range dbNames {
		connectDbServer(dbName)
	}
}

func connectDbServer(dbName string) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		utils.Configuration.Database.Host,
		utils.Configuration.Database.Port,
		utils.Configuration.Database.User,
		dbName,
	)
	if db, err := sql.Open("postgres", psqlInfo); err != nil {
		utils.LogError("Can't connect to database: %v", err)
	} else {
		ctx, stop := context.WithCancel(context.Background())
		defer stop()
		if pinged := ping(ctx, db); pinged {
			dbs[dbName] = db
			utils.Log("Database %q connected.", dbName)
		}
	}
}

func ping(ctx context.Context, db *sql.DB) bool {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		utils.LogError("Unable to connect to database: %v.", err)
		return false
	}
	return true
}

func DisconnectDbServers(dbNames []string) {
	for _, dbName := range dbNames {
		disconnectDbServer(dbName)
	}
}

func disconnectDbServer(dbName string) {
	if dbs[dbName] == nil {
		utils.Log("Skip database %q for disconnection.", dbName)
	} else {
		err := dbs[dbName].Close()
		if err != nil {
			utils.LogError("Unable to disconnect to database: %v.", err)
		} else {
			utils.Log("Database %q disconnected.", dbName)
		}
	}
}
