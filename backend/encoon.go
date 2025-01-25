// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"d.lambert.fr/encoon/apis"
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
)

var (
	configurationFileName string
	exportDb              string
	exportFileName        string
	exportAll             bool
	importDb              string
	importFileName        string
)

const (
	configurationFileNameFlag    = "configuration"
	defaultConfigurationFileName = "./configuration.yml"
	usageConfigurationFileName   = "Name of the file (.yml) used to configure the application (full path)"
	exportDbFlag                 = "export"
	defaultDbExport              = ""
	usageDbExport                = "Name of the database to export"
	exportFileNameFlag           = "exportfile"
	defaultExportFileName        = ""
	usageExportFileName          = "Name of the file (.ymp) used for exporting data"
	exportAllFlag                = "exportall"
	defaultExportAll             = false
	usageExportAll               = "Export all tables including including users and transactions"
	importDbFlag                 = "import"
	defaultDbImport              = ""
	usageDbImport                = "Name of the database to import"
	importFileNameFlag           = "importfile"
	defaultImportFileName        = ""
	usageImportFileName          = "Name of the file (.ymp) used for importing data"
)

func main() {
	handleFlags()
	if configuration.LoadConfiguration(configurationFileName) == nil {
		if exportDb != "" && exportFileName != "" {
			configuration.Log("", "", "Export")
			database.ExportDb(context.Background(), exportDb, exportFileName, exportAll)
		} else if importDb != "" && importFileName != "" {
			configuration.Log("", "", "Import")
			database.SeedDb(context.Background(), importDb, importFileName)
		} else {
			configuration.Log("", "", "Starting")
			configuration.WatchConfigurationChanges(configurationFileName)
			quitChan, doneChan := make(chan os.Signal, 1), make(chan bool, 1)
			signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-quitChan
				configuration.Log("", "", "Stopping")
				doneChan <- true
			}()
			go apis.InitializeCaches()
			go apis.StartReadingMessages()
			<-doneChan
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			apis.StopReadingMessages()
			<-ctx.Done()
		}
	}
	configuration.Log("", "", "Stopped")
}

func handleFlags() {
	flag.StringVar(&configurationFileName, configurationFileNameFlag, defaultConfigurationFileName, usageConfigurationFileName)
	flag.StringVar(&exportDb, exportDbFlag, defaultDbExport, usageDbExport)
	flag.StringVar(&exportFileName, exportFileNameFlag, defaultExportFileName, usageExportFileName)
	flag.BoolVar(&exportAll, exportAllFlag, defaultExportAll, usageExportAll)
	flag.StringVar(&importDb, importDbFlag, defaultDbImport, usageDbImport)
	flag.StringVar(&importFileName, importFileNameFlag, defaultImportFileName, usageImportFileName)
	flag.Parse()
}
