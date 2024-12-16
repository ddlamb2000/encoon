// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package main

import (
	"context"
	"flag"
	"fmt"
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
	exportDb              string
	exportFileName        string
)

const (
	configurationFileNameFlag    = "configuration"
	defaultConfigurationFileName = "./configuration.yml"
	usageConfigurationFileName   = "Name of the file (.yml) used to configure the application (full path)."
	exportDbFlag                 = "export"
	defaultDbExport              = ""
	usageDbExport                = "Name of the database to export."
	exportFileNameFlag           = "exportfile"
	defaultExportFileName        = ""
	usageExportFileName          = "Name of the file (.ymp) used for exporting data."
)

func main() {
	handleFlags()
	router.Use(gin.Logger())
	if configuration.LoadConfiguration(configurationFileName) == nil {
		if exportDb != "" && exportFileName != "" {
			database.ExportDb(context.Background(), exportDb, exportFileName)
		} else {
			configuration.Log("", "", "Starting...")
			configuration.WatchConfigurationChanges(configurationFileName)
			quitChan, doneChan := make(chan os.Signal), make(chan bool, 1)
			signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-quitChan
				configuration.Log("", "", "Stopping.")
				doneChan <- true
			}()
			go setAndStartHttpServer()
			go apis.InitializeCaches()
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
	flag.StringVar(&exportDb, exportDbFlag, defaultDbExport, usageDbExport)
	flag.StringVar(&exportFileName, exportFileNameFlag, defaultExportFileName, usageExportFileName)
	flag.Parse()
}

func setAndStartHttpServer() error {
	router.LoadHTMLGlob("frontend/templates/*.html")
	router.Static("/react", "./frontend/react")
	router.Static("/javascript", "./frontend/javascript")
	router.Static("/bootstrap", "./frontend/lib/bootstrap-5.2.3-dist")
	router.Static("/icons", "./frontend/lib/bootstrap-icons")
	router.Static("/quill", "./frontend/lib/quill")
	router.Static("/stylesheets", "./frontend")
	router.StaticFile("favicon.ico", "./frontend/favicon.ico")
	router.GET("/", getIndexHtml)
	router.GET("/:dbName", getIndexHtml)
	router.GET("/:dbName/:gridUuid", getIndexHtml)
	router.GET("/:dbName/:gridUuid/:uuid", getIndexHtml)
	apis.SetApiRoutes(router)
	httpPort := configuration.GetConfiguration().HttpServer.Port
	srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	configuration.Log("", "", "Listening http port %d.", httpPort)
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
