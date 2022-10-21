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
	middleware.SetAndStartHttpServer()
	err := middleware.SetAndStartDbServer()
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
	middleware.ShutDownHttpServer(ctx)
	middleware.ShutDownDbServer()
}
