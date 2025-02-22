// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"d.lambert.fr/encoon/configuration"
	_ "github.com/lib/pq"
)

var dbs = struct {
	sync.RWMutex
	m map[string]*sql.DB
}{m: make(map[string]*sql.DB)}

func GetDbByName(dbName string) (*sql.DB, error) {
	if dbName == "" || dbName == "undefined" {
		return nil, configuration.LogAndReturnError("", "", "Missing database name parameter.")
	}
	dbs.RLock()
	db := dbs.m[dbName]
	dbs.RUnlock()
	if db != nil {
		return db, nil
	}
	dbs.Lock()
	defer dbs.Unlock()
	dbConfiguration, err := findDbConfiguration(dbName)
	if err != nil {
		return nil, err
	}
	db, err = connectDbServer(dbConfiguration, dbName)
	if err != nil {
		return nil, err
	}
	dbs.m[dbName] = db
	return db, nil
}

func findDbConfiguration(dbName string) (*configuration.DatabaseConfiguration, error) {
	if configuration.GetConfiguration().Databases != nil {
		for _, conf := range configuration.GetConfiguration().Databases {
			if conf.Name == dbName {
				return conf, nil
			}
		}
	}
	return nil, configuration.LogAndReturnError(dbName, "", "Database isn't configured.")
}

func connectDbServer(dbConfiguration *configuration.DatabaseConfiguration, dbName string) (*sql.DB, error) {
	if dbConfiguration.Host == "" || dbConfiguration.Port == 0 || dbConfiguration.Role == "" || dbConfiguration.Name == "" {
		return nil, configuration.LogAndReturnError(dbName, "", "Incorrect database configuration.")
	}
	var dbHost = dbConfiguration.Host
	mockDbHost := os.Getenv("ENCOON_MOCK_DB_LOCALHOST")
	if dbHost == "localhost" && mockDbHost != "" {
		dbHost = mockDbHost
		configuration.Log(dbName, "", "Use mock host %s.", mockDbHost)
	}
	configuration.Log(dbName,
		"",
		"Connect to database host=%s port=%d user=%s.",
		dbHost,
		dbConfiguration.Port,
		dbConfiguration.Role)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		dbHost,
		dbConfiguration.Port,
		dbConfiguration.Role,
		dbConfiguration.RolePassword,
		dbConfiguration.Name,
	)
	db, _ := sql.Open("postgres", psqlInfo)
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := db.PingContext(ctx); err != nil {
		return nil, configuration.LogAndReturnError(dbName, "", "Unable to connect to database: %v.", err)
	}
	configuration.Log(dbName, "", "Database connected.")
	return db, migrateDb(ctx, db, dbConfiguration.Name)
}

func Sleep(ctx context.Context, dbName, user string, db *sql.DB) {
	dbConfiguration := configuration.GetDatabaseConfiguration(dbName)
	sleepTime := dbConfiguration.TestSleepTime
	if sleepTime > 0 {
		configuration.Trace(dbName, user, "*** START SLEEP *** It's now %v.", time.Now())
		wait := float32(sleepTime) * (1.0 + rand.Float32()) / 2.0 / 1000.0
		st := fmt.Sprintf("SELECT pg_sleep(%.2f)", wait)
		db.QueryContext(ctx, st)
		configuration.Trace(dbName, user, "*** STOP SLEEP *** It's now %v.", time.Now())
	}
}

func ForceTestSleepTimeAndTimeOutThreshold(dbName string, sleepTime, threshold int) {
	dbConfiguration := configuration.GetDatabaseConfiguration(dbName)
	if dbConfiguration != nil {
		dbConfiguration.TestSleepTime = sleepTime
		dbConfiguration.TimeOutThreshold = threshold
	}
}
