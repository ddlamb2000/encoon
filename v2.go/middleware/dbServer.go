// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"d.lambert.fr/encoon/backend/utils"
	"golang.org/x/exp/maps"
)

var (
	dbs = make(map[string]*sql.DB)
)

func ConnectDbServers(dbConfigurations map[string]*utils.DatabaseConfig) {
	for _, conf := range maps.Values(dbConfigurations) {
		connectDbServer(conf)
	}
}

func connectDbServer(dbConfiguration *utils.DatabaseConfig) {
	dbName := dbConfiguration.Database.Name
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		dbConfiguration.Database.Host,
		dbConfiguration.Database.Port,
		dbConfiguration.Database.User,
		dbConfiguration.Database.Name,
	)
	if db, err := sql.Open("postgres", psqlInfo); err != nil {
		utils.LogError("Can't connect to database: %v", err)
	} else {
		ctx, stop := context.WithCancel(context.Background())
		defer stop()
		if pinged := ping(ctx, db); pinged {
			dbs[dbName] = db
			utils.Log("Database %q connected.", dbName)
			migrateDb(ctx, db)
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

func DisconnectDbServers(dbConfigurations map[string]*utils.DatabaseConfig) {
	for _, conf := range maps.Values(dbConfigurations) {
		disconnectDbServer(conf)
	}
}

func disconnectDbServer(dbConfiguration *utils.DatabaseConfig) {
	dbName := dbConfiguration.Database.Name
	if dbs[dbName] != nil {
		err := dbs[dbName].Close()
		if err != nil {
			utils.LogError("Unable to disconnect to database: %v.", err)
		} else {
			utils.Log("Database %q disconnected.", dbName)
		}
	}
}
