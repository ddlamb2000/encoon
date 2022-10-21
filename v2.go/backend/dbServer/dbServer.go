// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package dbServer

import (
	"database/sql"
	"fmt"

	"d.lambert.fr/encoon/backend/utils"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "david.lambert"
	dbPassword = ""
	dbName     = "david.lambert"
)

var (
	db  *sql.DB
	err error
)

func SetAndStartServer() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	utils.Logf("Database connected on port %d", dbPort)
	return err
}

func ShutDownServer() {
	utils.Log("Database disconnection...")
	db.Close()
	utils.Log("Database closed.")
}
