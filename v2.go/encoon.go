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
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router.Use(gin.Logger())
	if configuration.LoadConfiguration(configurationFileName) == nil {
		configuration.Log("Starting, log into %v.", logFileName)
		if exportDb != "" && exportFileName != "" {
			database.ExportDb(context.Background(), exportDb, exportFileName)
			configuration.Log("Exported database %s into %s.", exportDb, exportFileName)
		} else {
			configuration.WatchConfigurationChanges(configurationFileName)
			quitChan, doneChan := make(chan os.Signal), make(chan bool, 1)
			signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-quitChan
				configuration.Log("Stopping.")
				doneChan <- true
			}()
			go setAndStartHttpServer()
			<-doneChan
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				configuration.LogError("Error during server shutdown: %v.", err)
			}
			select {
			case <-ctx.Done():
				configuration.Log("Timeout of 5 seconds.")
			}
		}
	}
	configuration.Log("Stopped.")
}

func handleFlags() {
	flag.StringVar(&configurationFileName, configurationFileNameFlag, defaultConfigurationFileName, usageConfigurationFileName)
	flag.StringVar(&logFileName, logFileNameFlag, defaultLogFileName, usageLogFileName)
	flag.StringVar(&exportDb, exportDbFlag, defaultDbExport, usageDbExport)
	flag.StringVar(&exportFileName, exportFileNameFlag, defaultExportFileName, usageExportFileName)
	flag.Parse()
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
	apis.SetApiRoutes(router)
	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", configuration.GetConfiguration().HttpServer.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	configuration.Log("Listening http.")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		configuration.LogError("Error on http listening: %v.", err)
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
