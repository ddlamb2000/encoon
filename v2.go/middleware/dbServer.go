// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"database/sql"
	"fmt"

	"d.lambert.fr/encoon/backend/utils"
)

const (
	_dbHost     = "localhost"
	_dbPort     = 5432
	_dbUser     = "david.lambert"
	_dbPassword = ""
	_dbName     = "david.lambert"
)

var (
	db  *sql.DB
	err error
)

func SetAndStartDbServer() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", _dbHost, _dbPort, _dbUser, _dbName)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	utils.Logf("Database connected on port %d", _dbPort)
}

func ShutDownDbServer() {
	utils.Log("Database disconnection...")
	err = db.Close()
	if err != nil {
		panic(err)
	}
	utils.Log("Database closed.")
}
