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
	"d.lambert.fr/encoon/utils"
)

func main() {
	utils.InitWithLog()
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		<-quitChan
		utils.Log("Stopping.")
		done <- true
	}()
	if utils.LoadConfiguration("configurations/") == nil {
		go backend.ConnectDbServers(utils.DatabaseConfigurations)
		go backend.SetAndStartHttpServer()
		<-done
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		backend.ShutDownHttpServer(ctx)
		backend.DisconnectDbServers(utils.DatabaseConfigurations)
		select {
		case <-ctx.Done():
			utils.Log("Timeout of 5 seconds.")
		}
	}
	utils.Log("Stopped.")
}
