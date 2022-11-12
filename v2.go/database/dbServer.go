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

var (
	dbs = make(map[string]*sql.DB)
)

func GetDbByName(dbName string) *sql.DB {
	return dbs[dbName]
}

func setDb(dbName string, db *sql.DB) {
	dbs[dbName] = db
}

func ConnectDbServers(dbConfigurations []*configuration.DatabaseConfiguration) error {
	for _, conf := range dbConfigurations {
		if err := connectDbServer(conf); err != nil {
			return err
		}
	}
	return nil
}

func connectDbServer(dbConfiguration *configuration.DatabaseConfiguration) error {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		dbConfiguration.Host,
		dbConfiguration.Port,
		dbConfiguration.User,
		dbConfiguration.Name,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		utils.LogError("Can't connect to database with %q: %v", psqlInfo, err)
		return err
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := PingDb(ctx, db); err != nil {
		return err
	}
	setDb(dbConfiguration.Name, db)
	utils.Log("[%s] Database connected.", dbConfiguration.Name)
	migrateDb(ctx, db, dbConfiguration.Name)
	return nil
}

func PingDb(ctx context.Context, db *sql.DB) error {
	if err := db.PingContext(context.Background()); err != nil {
		utils.LogError("Unable to connect to database: %v.", err)
		return err
	}
	return nil
}

func DisconnectDbServers(dbConfigurations []*configuration.DatabaseConfiguration) error {
	for _, conf := range dbConfigurations {
		if err := disconnectDbServer(conf); err != nil {
			return err
		}
	}
	return nil
}

func disconnectDbServer(dbConfiguration *configuration.DatabaseConfiguration) error {
	dbName := dbConfiguration.Name
	db := GetDbByName(dbName)
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
