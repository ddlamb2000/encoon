// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"database/sql"
	"fmt"

	"d.lambert.fr/encoon/backend/utils"
)

var (
	db  *sql.DB
	err error
)

func SetAndStartDbServer() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		utils.Configuration.Database.Host,
		utils.Configuration.Database.Port,
		utils.Configuration.Database.User,
		utils.Configuration.Database.Name,
	)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	utils.Log(fmt.Sprintf("Database %s connected.", utils.Configuration.Database.Name))
}

func ShutDownDbServer() {
	utils.Log(fmt.Sprintf("Database %s disconnection.", utils.Configuration.Database.Name))
	err = db.Close()
	if err != nil {
		panic(err)
	}
	utils.Log(fmt.Sprintf("Database %s disconnected.", utils.Configuration.Database.Name))
}
