// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"d.lambert.fr/encoon/utils"
	_ "github.com/lib/pq"
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

func ConnectDbServers(dbConfigurations map[string]*utils.DatabaseConfig) error {
	for _, conf := range maps.Values(dbConfigurations) {
		if err := connectDbServer(conf); err != nil {
			return err
		}
	}
	return nil
}

func connectDbServer(dbConfiguration *utils.DatabaseConfig) error {
	dbName := dbConfiguration.Database.Name
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		dbConfiguration.Database.Host,
		dbConfiguration.Database.Port,
		dbConfiguration.Database.User,
		dbConfiguration.Database.Name,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		utils.LogError("Can't connect to database with %q: %v", psqlInfo, err)
		return err
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := pingDb(ctx, db); err != nil {
		return err
	}
	setDb(dbName, db)
	utils.Log("[%s] Database connected.", dbName)
	migrateDb(ctx, db, dbName)
	return nil
}

func pingDb(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		utils.LogError("Unable to connect to database: %v.", err)
		return err
	}
	return nil
}

func DisconnectDbServers(dbConfigurations map[string]*utils.DatabaseConfig) error {
	for _, conf := range maps.Values(dbConfigurations) {
		if err := disconnectDbServer(conf); err != nil {
			return err
		}
	}
	return nil
}

func disconnectDbServer(dbConfiguration *utils.DatabaseConfig) error {
	dbName := dbConfiguration.Database.Name
	db := getDbByName(dbName)
	if db != nil {
		err := db.Close()
		if err != nil {
			utils.LogError("Unable to disconnect database %q: %v.", dbName, err)
			return err
		}
		utils.Log("[%s] Database disconnected.", dbName)
	}
	return nil
}

func beginTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid string, user string) error {
	_, err := db.ExecContext(ctx, " BEGIN ")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Begin transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Begin transaction.", dbName, user)
	return err
}

func commitTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid string, user string) error {
	_, err := db.ExecContext(ctx, " COMMIT ")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Commit transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Commit transaction.", dbName, user)
	return err
}
