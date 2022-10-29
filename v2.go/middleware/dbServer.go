// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"d.lambert.fr/encoon/utils"
	"golang.org/x/exp/maps"
)

var (
	dbs = make(map[string]*sql.DB)
)

func getDbByName(dbName string) *sql.DB {
	return dbs[dbName]
}

func setDb(dbName string, db *sql.DB) {
	dbs[dbName] = db
}

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
		if pinged := pingDb(ctx, db); pinged {
			setDb(dbName, db)
			utils.Log("Database %q connected.", dbName)
			migrateDb(ctx, db, dbName)
		}
	}
}

func pingDb(ctx context.Context, db *sql.DB) bool {
	if db == nil {
		utils.LogError("No database provided for ping.")
		return false
	}
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
	db := getDbByName(dbName)
	if db != nil {
		err := db.Close()
		if err != nil {
			utils.LogError("Unable to disconnect database %q: %v.", dbName, err)
		} else {
			utils.Log("Database %q disconnected.", dbName)
		}
	}
}
