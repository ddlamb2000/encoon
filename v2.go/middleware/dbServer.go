// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"database/sql"
	"fmt"

	"d.lambert.fr/encoon/backend/utils"
)

var dbs = make(map[string]*sql.DB)

func ConnectDbServers() {
	for _, v := range utils.Configuration.DatabaseNames {
		psqlInfo := fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable",
			utils.Configuration.Database.Host,
			utils.Configuration.Database.Port,
			utils.Configuration.Database.User,
			v,
		)
		var err error
		dbs[v], err = sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		utils.Log(fmt.Sprintf("Database %s connected.", v))
	}
}

func DisconnectDbServers() {
	for _, v := range utils.Configuration.DatabaseNames {
		err := dbs[v].Close()
		if err != nil {
			panic(err)
		}
		utils.Log(fmt.Sprintf("Database %s disconnected.", v))
	}
}
