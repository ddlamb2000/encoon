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

	"d.lambert.fr/encoon/backend"
	"d.lambert.fr/encoon/middleware"
	"d.lambert.fr/encoon/utils"
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

	utils.LoadConfiguration("configurations/")
	go middleware.ConnectDbServers(utils.DatabaseConfigurations)
	go middleware.SetAndStartHttpServer()
	go backend.LoadData()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	middleware.ShutDownHttpServer(ctx)
	middleware.DisconnectDbServers(utils.DatabaseConfigurations)
	utils.Log("Stopped.")
}
