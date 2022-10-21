// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"d.lambert.fr/encoon/backend/core"
	"d.lambert.fr/encoon/backend/dbServer"
	"d.lambert.fr/encoon/backend/httpServer"
	"d.lambert.fr/encoon/backend/utils"
)

func main() {
	utils.InitWithLog()
	httpServer.SetAndStartServer()
	err := dbServer.SetAndStartServer()
	if err != nil {
		panic(err)
	}
	core.LoadData()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log("Shut down (SIGTERM)...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.ShutDownServer(ctx)
	dbServer.ShutDownServer()
}
