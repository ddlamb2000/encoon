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
	"d.lambert.fr/encoon/backend/utils"
	"d.lambert.fr/encoon/middleware"
)

func main() {
	utils.InitWithLog()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-quit
		utils.Log("Stopping.")
		done <- true
	}()

	utils.LoadConfiguration()
	go middleware.ConnectDbServers(utils.Configuration.DatabaseNames)
	go middleware.SetAndStartHttpServer()
	go core.LoadData()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	middleware.ShutDownHttpServer(ctx)
	middleware.DisconnectDbServers(utils.Configuration.DatabaseNames)
	utils.Log("Stopped.")
}
