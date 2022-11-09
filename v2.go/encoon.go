// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"d.lambert.fr/encoon/backend"
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
	if utils.LoadConfiguration("configurations/") == nil {
		go backend.ConnectDbServers(utils.DatabaseConfigurations)
		go backend.SetAndStartHttpServer()
		<-done
		ctx, cancel := utils.GetContextWithTimeOut()
		defer cancel()
		backend.ShutDownHttpServer(ctx)
		backend.DisconnectDbServers(utils.DatabaseConfigurations)
	}
	utils.Log("Stopped.")
}
