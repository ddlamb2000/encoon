// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"d.lambert.fr/encoon/apis"
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
)

var (
	router                = gin.New()
	srv                   *http.Server
	configurationFileName string
	logFileName           string
	exportDb              string
	exportFileName        string
)

const (
	configurationFileNameFlag    = "configuration"
	defaultConfigurationFileName = "./configuration.yml"
	usageConfigurationFileName   = "Name of the file (.yml) used to configure the application (full path)."
	logFileNameFlag              = "log"
	defaultLogFileName           = "./logs/encoon.log"
	usageLogFileName             = "Name of the file (.log) used for logging."
	exportDbFlag                 = "export"
	defaultDbExport              = ""
	usageDbExport                = "Name of the database to export."
	exportFileNameFlag           = "exportfile"
	defaultExportFileName        = ""
	usageExportFileName          = "Name of the file (.ymp) used for exporting data."
)

func main() {
	handleFlags()
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(logFileName, flags, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router.Use(gin.Logger())
	if configuration.LoadConfiguration(configurationFileName) == nil {
		if exportDb != "" && exportFileName != "" {
			database.ExportDb(context.Background(), exportDb, exportFileName)
		} else {
			configuration.Log("", "", "Starting, log into %v.", logFileName)
			configuration.WatchConfigurationChanges(configurationFileName)
			quitChan, doneChan := make(chan os.Signal), make(chan bool, 1)
			signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-quitChan
				configuration.Log("", "", "Stopping.")
				doneChan <- true
			}()
			go apis.InitializeCaches()
			go setAndStartHttpServer()
			<-doneChan
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				configuration.LogError("", "", "Error during server shutdown: %v.", err)
			}
			<-ctx.Done()
		}
	}
	configuration.Log("", "", "Stopped.")
}

func handleFlags() {
	flag.StringVar(&configurationFileName, configurationFileNameFlag, defaultConfigurationFileName, usageConfigurationFileName)
	flag.StringVar(&logFileName, logFileNameFlag, defaultLogFileName, usageLogFileName)
	flag.StringVar(&exportDb, exportDbFlag, defaultDbExport, usageDbExport)
	flag.StringVar(&exportFileName, exportFileNameFlag, defaultExportFileName, usageExportFileName)
	flag.Parse()
}

func setAndStartHttpServer() error {
	if configuration.IsFrontEndDevelopment() {
		router.LoadHTMLGlob("frontend/templates/index-development.html")
		router.Static("/javascript", "./frontend/react")
	} else {
		router.LoadHTMLGlob("frontend/templates/index-production.html")
		router.Static("/javascript", "./frontend/javascript")
	}
	router.Static("/bootstrap", "./frontend/bootstrap-5.2.3-dist")
	router.Static("/stylesheets", "./frontend/stylesheets")
	router.Static("/images", "./frontend/images")
	router.Static("/icons", "./frontend/bootstrap-icons/icons")
	router.StaticFile("favicon.ico", "./frontend/images/favicon.ico")
	router.GET("/", getIndexHtml)
	router.GET("/:dbName", getIndexHtml)
	router.GET("/:dbName/:gridUuid", getIndexHtml)
	router.GET("/:dbName/:gridUuid/:uuid", getIndexHtml)
	apis.SetApiRoutes(router)
	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", configuration.GetConfiguration().HttpServer.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	configuration.Log("", "", "Listening http.")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		configuration.LogError("", "", "Error on http listening: %v.", err)
		return err
	}
	return nil
}

func getIndexHtml(c *gin.Context) {
	var indexFile = "index-production.html"
	if configuration.IsFrontEndDevelopment() {
		indexFile = "index-development.html"
	}
	c.HTML(http.StatusOK, indexFile, gin.H{
		"appName":  configuration.GetConfiguration().AppName,
		"appTag":   configuration.GetConfiguration().AppTag,
		"dbName":   c.Param("dbName"),
		"gridUuid": c.Param("gridUuid"),
		"uuid":     c.Param("uuid"),
	})
}
