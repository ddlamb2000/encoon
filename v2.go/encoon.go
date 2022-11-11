// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"d.lambert.fr/encoon/backend"
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
)

var (
	router = gin.Default()
	srv    *http.Server
)

func main() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	if configuration.LoadConfiguration("./", "configuration.yml") == nil {
		utils.Log("Starting.")
		quitChan, doneChan := make(chan os.Signal), make(chan bool, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-quitChan
			utils.Log("Stopping.")
			doneChan <- true
		}()
		go database.ConnectDbServers(configuration.GetConfiguration().Databases)
		go setAndStartHttpServer()
		<-doneChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go database.DisconnectDbServers(configuration.GetConfiguration().Databases)
		if err := srv.Shutdown(ctx); err != nil {
			utils.LogError("Error during server shutdown: %v.", err)
		}
		select {
		case <-ctx.Done():
			utils.Log("Timeout of 5 seconds.")
		}
	}
	utils.Log("Stopped.")
}

func setAndStartHttpServer() error {
	router.LoadHTMLGlob("frontend/templates/index.html")
	router.Static("/stylesheets", "./frontend/stylesheets")
	router.Static("/javascript", "./frontend/javascript")
	router.Static("/images", "./frontend/images")
	router.Static("/icons", "./frontend/bootstrap-icons/icons")
	router.StaticFile("favicon.ico", "./frontend/images/favicon.ico")
	router.GET("/", getIndexHtml)
	router.GET("/:dbName", getIndexHtml)
	router.GET("/:dbName/:gridUri", getIndexHtml)
	router.GET("/:dbName/:gridUri/:uuid", getIndexHtml)
	backend.SetApiRoutes(router)
	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", configuration.GetConfiguration().HttpServer.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	utils.Log("Listening http.")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.LogError("Error on http listening: %v.", err)
		return err
	}
	return nil
}

func getIndexHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"appName": configuration.GetConfiguration().AppName,
		"appTag":  configuration.GetConfiguration().AppTag,
		"dbName":  c.Param("dbName"),
		"gridUri": c.Param("gridUri"),
		"uuid":    c.Param("uuid"),
	})
}
