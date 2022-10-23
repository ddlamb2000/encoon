// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			utils.LogFatal("Database:", err)
		}

		ctx, stop := context.WithCancel(context.Background())
		defer stop()

		ping(ctx, db)
		dbs[v] = db
		utils.Log(fmt.Sprintf("Database %s connected.", v))
	}
}

func ping(ctx context.Context, db *sql.DB) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		utils.LogFatal(fmt.Sprintf("Unable to connect to database: %v.", err))
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
