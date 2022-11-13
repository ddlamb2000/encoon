// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/utils"
	_ "github.com/lib/pq"
)

var dbs = make(map[string]*sql.DB)

func GetDbByName(dbName string) (*sql.DB, error) {
	if dbName == "" {
		return nil, utils.LogAndReturnError("Missing database name parameter.")
	}
	db := dbs[dbName]
	if db != nil {
		return db, nil
	}
	dbConfiguration, err := findDbConfiguration(dbName)
	if err != nil {
		return nil, err
	}
	db, err = connectDbServer(dbConfiguration)
	if err != nil {
		return nil, err
	}
	dbs[dbName] = db
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
	return nil, utils.LogAndReturnError("Database %q isn't configured.", dbName)
}

func connectDbServer(dbConfiguration *configuration.DatabaseConfiguration) (*sql.DB, error) {
	if dbConfiguration.Host == "" || dbConfiguration.Port == 0 || dbConfiguration.User == "" || dbConfiguration.Name == "" {
		return nil, utils.LogAndReturnError("Incorrect database configuration.")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		dbConfiguration.Host,
		dbConfiguration.Port,
		dbConfiguration.User,
		dbConfiguration.Name,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, utils.LogAndReturnError("Can't connect to database with %q: %v", psqlInfo, err)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := PingDb(ctx, db); err != nil {
		return nil, err
	}
	utils.Log("[%s] Database connected.", dbConfiguration.Name)
	migrateDb(ctx, db, dbConfiguration.Name)
	return db, nil
}

func PingDb(ctx context.Context, db *sql.DB) error {
	if err := db.PingContext(ctx); err != nil {
		return utils.LogAndReturnError("Unable to connect to database: %v.", err)
	}
	return nil
}

func BeginTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user, trace string) error {
	utils.Trace(trace, "beginTransaction()")
	_, err := db.ExecContext(ctx, "BEGIN")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Begin transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Begin transaction.", dbName, user)
	return err
}

func CommitTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user, trace string) error {
	utils.Trace(trace, "commitTransaction()")
	_, err := db.ExecContext(ctx, "COMMIT")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Commit transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] Commit transaction.", dbName, user)
	return err
}

func RollbackTransaction(ctx context.Context, dbName string, db *sql.DB, userUuid, user, trace string) error {
	utils.Trace(trace, "rollbackTransaction()")
	_, err := db.ExecContext(ctx, "ROLLBACK")
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Rollback transaction error: %v.", dbName, user, err)
	}
	utils.Log("[%s] [%s] ROLLBACK transaction.", dbName, user)
	return err
}

func TestSleep(ctx context.Context, dbName string, db *sql.DB) error {
	dbConfiguration := configuration.GetDatabaseConfiguration(dbName)
	sleepTime := dbConfiguration.TestSleepTime
	if sleepTime > 0 {
		wait := float32(sleepTime) * (1.0 + rand.Float32()) / 2.0 / 1000.0
		st := fmt.Sprintf("SELECT pg_sleep(%.2f)", wait)
		utils.Log("*** BEFORE SLEEP It's now %v.", time.Now())
		_, err := db.QueryContext(ctx, st)
		utils.Log("*** AFTER SLEEP It's now %v.", time.Now())
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
