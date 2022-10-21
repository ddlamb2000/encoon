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
	middleware.SetAndStartDbServer()
	go middleware.SetAndStartHttpServer()
	go core.LoadData()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log("Stopping.")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	middleware.ShutDownHttpServer(ctx)
	middleware.ShutDownDbServer()
	select {
	case <-ctx.Done():
		utils.Log("Stopped.")
	}
}
