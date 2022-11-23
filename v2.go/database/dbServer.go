// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
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
	if dbName == "" {
		return nil, configuration.LogAndReturnError("", "", "Missing database name parameter.")
	}
	dbs.RLock()
	db := dbs.m[dbName]
	dbs.RUnlock()
	if db != nil {
		return db, nil
	}
	dbConfiguration, err := findDbConfiguration(dbName)
	if err != nil {
		return nil, err
	}
	db, err = connectDbServer(dbConfiguration, dbName)
	if err != nil {
		return nil, err
	}
	dbs.Lock()
	dbs.m[dbName] = db
	dbs.Unlock()
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
	if dbConfiguration.Host == "" || dbConfiguration.Port == 0 || dbConfiguration.User == "" || dbConfiguration.Name == "" {
		return nil, configuration.LogAndReturnError(dbName, "", "Incorrect database configuration.")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		dbConfiguration.Host,
		dbConfiguration.Port,
		dbConfiguration.User,
		dbConfiguration.Name,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, configuration.LogAndReturnError(dbName, "", "Can't connect to database with %q: %v", psqlInfo, err)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := db.PingContext(ctx); err != nil {
		return nil, configuration.LogAndReturnError(dbName, "", "Unable to connect to database: %v.", err)
	}
	configuration.Log(dbName, "", "Database connected.")
	migrateDb(ctx, db, dbConfiguration.Name)
	return db, nil
}

func TestSleep(ctx context.Context, dbName, user string, db *sql.DB) error {
	dbConfiguration := configuration.GetDatabaseConfiguration(dbName)
	sleepTime := dbConfiguration.TestSleepTime
	if sleepTime > 0 {
		wait := float32(sleepTime) * (1.0 + rand.Float32()) / 2.0 / 1000.0
		st := fmt.Sprintf("SELECT pg_sleep(%.2f)", wait)
		configuration.Log(dbName, user, "*** BEFORE SLEEP It's now %v.", time.Now())
		_, err := db.QueryContext(ctx, st)
		configuration.Log(dbName, user, "*** AFTER SLEEP It's now %v.", time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func ForceTestSleepTimeAndTimeOutThreshold(dbName string, sleepTime, threshold int) {
	dbConfiguration := configuration.GetDatabaseConfiguration(dbName)
	if dbConfiguration != nil {
		dbConfiguration.TestSleepTime = sleepTime
		dbConfiguration.TimeOutThreshold = threshold
	}
}
