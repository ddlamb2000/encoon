// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package configuration

import (
	"context"
	"os"
	"sync"
	"time"

	"d.lambert.fr/encoon/utils"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	AppName string `yaml:"appName"`
	AppTag  string `yaml:"appTag"`

	HttpServer struct {
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		JwtExpiration int    `yaml:"jwtExpiration"`
	} `yaml:"httpServer"`

	Databases []*Database `yaml:"database"`

	valid bool
}

type Database struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	Name             string `yaml:"name"`
	User             string `yaml:"user"`
	JwtSecret        string `yaml:"jwtsecret"`
	Root             string `yaml:"root"`
	Password         string `yaml:"password"`
	TestSleepTime    int    `yaml:"testSleepTime"`
	TimeOutThreshold int    `yaml:"timeOutThreshold"`

	valid bool
}

var (
	appConfiguration      Configuration
	configurationFileName string
	appConfigurationMutex sync.Mutex
)

func LoadConfiguration(fileName string) error {
	configurationFileName = fileName
	return loadConfigurationFromFile()
}

func loadConfigurationFromFile() error {
	utils.Log("Loading configuration from %v.", configurationFileName)
	appConfigurationMutex.Lock()
	defer appConfigurationMutex.Unlock()
	f, err := os.Open(configurationFileName)
	if err != nil {
		return utils.LogAndReturnError("Error loading configuration from file %q: %v.", configurationFileName, err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	newConfiguration := new(Configuration)
	if err = decoder.Decode(&newConfiguration); err != nil {
		return utils.LogAndReturnError("Error parsing configuration from file %q: %v.", configurationFileName, err)
	}
	if err = validateConfiguration(newConfiguration); err != nil {
		return err
	}
	appConfiguration = *newConfiguration
	utils.Log("Configuration loaded from file %q.", configurationFileName)
	return nil
}

func validateConfiguration(conf *Configuration) error {
	if conf.AppName == "" {
		return utils.LogAndReturnError("Missing application name (appName) from configuration file %v.", configurationFileName)
	}
	if conf.AppTag == "" {
		return utils.LogAndReturnError("Missing application tag line (appTag) from configuration file %v.", configurationFileName)
	}
	if conf.HttpServer.Host == "" {
		return utils.LogAndReturnError("Missing host name (httpServer.host) from configuration file %v.", configurationFileName)
	}
	if conf.HttpServer.Port == 0 {
		return utils.LogAndReturnError("Missing port (httpServer.port) from configuration file %v.", configurationFileName)
	}
	if conf.HttpServer.JwtExpiration == 0 {
		return utils.LogAndReturnError("Missing expiration (httpServer.jwtExpiration) from configuration file %v.", configurationFileName)
	}
	conf.valid = true
	utils.Log("Configuration from %v is valid.", configurationFileName)
	return nil
}

func GetConfiguration() Configuration {
	appConfigurationMutex.Lock()
	defer appConfigurationMutex.Unlock()
	return appConfiguration
}

func IsConfigurationValid() bool {
	appConfigurationMutex.Lock()
	defer appConfigurationMutex.Unlock()
	return appConfiguration.valid
}

func IsDatabaseEnabled(dbName string) bool {
	return dbName != "" && GetDatabaseConfiguration(dbName) != nil
}

func GetDatabaseConfiguration(dbName string) *Database {
	for _, dbConfig := range appConfiguration.Databases {
		if dbConfig.Name == dbName {
			return dbConfig
		}
	}
	return nil
}

func GetJWTSecret(dbName string) []byte {
	if !IsDatabaseEnabled(dbName) {
		return nil
	}
	return []byte(dbName + GetDatabaseConfiguration(dbName).JwtSecret)
}

func GetRootAndPassword(dbName string) (string, string) {
	if !IsDatabaseEnabled(dbName) {
		return "", ""
	}
	dbConfiguration := GetDatabaseConfiguration(dbName)
	return dbConfiguration.Root, dbConfiguration.Password
}

func GetContextWithTimeOut(dbName string) (context.Context, context.CancelFunc) {
	var threshold = 0
	dbConfiguration := GetDatabaseConfiguration(dbName)
	if dbConfiguration != nil {
		threshold = dbConfiguration.TimeOutThreshold
	}
	if threshold < 10 {
		threshold = 10
	}
	ctx, ctxFunc := context.WithTimeout(context.Background(), time.Duration(threshold)*time.Millisecond)
	return ctx, ctxFunc
}
