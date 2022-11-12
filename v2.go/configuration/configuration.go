// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package configuration

import (
	"context"
	"io/ioutil"
	"sync"
	"time"

	"d.lambert.fr/encoon/utils"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	AppName    string                   `yaml:"appName"`
	AppTag     string                   `yaml:"appTag"`
	HttpServer HttpServerConfiguration  `yaml:"httpServer"`
	Databases  []*DatabaseConfiguration `yaml:"database"`

	valid bool
}

type HttpServerConfiguration struct {
	Port          int `yaml:"port"`
	JwtExpiration int `yaml:"jwtExpiration"`
}

type DatabaseConfiguration struct {
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
	configurationHash     string
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
	f, err := ioutil.ReadFile(configurationFileName)
	if err != nil {
		return utils.LogAndReturnError("Error loading configuration from file %q: %v.", configurationFileName, err)
	}
	newConfiguration := new(Configuration)
	if err = yaml.Unmarshal(f, &newConfiguration); err != nil {
		return utils.LogAndReturnError("Error parsing configuration from file %q: %v.", configurationFileName, err)
	}
	if err = validateConfiguration(newConfiguration); err != nil {
		return err
	}
	hash, err := utils.CalculateFileHash(configurationFileName)
	if err != nil {
		return utils.LogAndReturnError("Error when calculating hash for configuration file %q: %v.", configurationFileName, err)
	}
	appConfiguration = *newConfiguration
	configurationHash = hash
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

func GetDatabaseConfiguration(dbName string) *DatabaseConfiguration {
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

func WatchConfigurationChanges(fileName string) {
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			newHash, err := utils.CalculateFileHash(fileName)
			if err != nil {
				utils.LogError("Error watching configuration changes on file %q: %v.", fileName, err)
				continue
			}
			if newHash != configurationHash {
				configurationHash = newHash
				LoadConfiguration(fileName)
			}
		}
	}()
}
