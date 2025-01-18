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

	if configuration.LoadConfiguration(configurationFileName) == nil {
		if exportDb != "" && exportFileName != "" {
			database.ExportDb(context.Background(), exportDb, exportFileName)
		} else {
			configuration.Log("", "", "Starting...")
			configuration.WatchConfigurationChanges(configurationFileName)
			quitChan, doneChan := make(chan os.Signal, 1), make(chan bool, 1)
			signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-quitChan
				configuration.Log("", "", "Stopping...")
				doneChan <- true
			}()
			go apis.InitializeCaches()
			go apis.ReadMessagesFromKafka()
			<-doneChan
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			apis.ShutdownKafkaProducers()
			apis.ShutdownKafkaConsumers()
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
