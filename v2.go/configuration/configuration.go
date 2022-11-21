// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package configuration

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	AppName    string                   `yaml:"appName"`
	AppTag     string                   `yaml:"appTag"`
	Trace      bool                     `yaml:"trace"`
	HttpServer HttpServerConfiguration  `yaml:"httpServer"`
	Databases  []*DatabaseConfiguration `yaml:"database"`
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
	Log("", "", "Loading configuration from %v.", configurationFileName)
	appConfigurationMutex.Lock()
	defer appConfigurationMutex.Unlock()
	f, err := ioutil.ReadFile(configurationFileName)
	if err != nil {
		return LogAndReturnError("", "", "Error loading configuration from file %q: %v.", configurationFileName, err)
	}
	newConfiguration := new(Configuration)
	if err = yaml.Unmarshal(f, &newConfiguration); err != nil {
		return LogAndReturnError("", "", "Error parsing configuration from file %q: %v.", configurationFileName, err)
	}
	if err = validateConfiguration(newConfiguration); err != nil {
		return err
	}
	hash, err := utils.CalculateFileHash(configurationFileName)
	if err != nil {
		return err
	}
	appConfiguration = *newConfiguration
	configurationHash = hash
	Log("", "", "Configuration loaded from file %q.", configurationFileName)
	return nil
}

func validateConfiguration(conf *Configuration) error {
	if conf.AppName == "" {
		return LogAndReturnError("", "", "Missing application name (appName) from configuration file %v.", configurationFileName)
	}
	if conf.AppTag == "" {
		return LogAndReturnError("", "", "Missing application tag line (appTag) from configuration file %v.", configurationFileName)
	}
	if conf.HttpServer.Port == 0 {
		return LogAndReturnError("", "", "Missing port (httpServer.port) from configuration file %v.", configurationFileName)
	}
	if conf.HttpServer.JwtExpiration == 0 {
		return LogAndReturnError("", "", "Missing expiration (httpServer.jwtExpiration) from configuration file %v.", configurationFileName)
	}
	Log("", "", "Configuration from %v is valid.", configurationFileName)
	return nil
}

func GetConfiguration() Configuration {
	return appConfiguration
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

func GetContextWithTimeOut(ct context.Context, dbName string) (context.Context, context.CancelFunc) {
	var threshold = 0
	dbConfiguration := GetDatabaseConfiguration(dbName)
	if dbConfiguration != nil {
		threshold = dbConfiguration.TimeOutThreshold
	}
	if threshold < 10 {
		threshold = 10
	}
	ctx, ctxFunc := context.WithTimeout(ct, time.Duration(threshold)*time.Millisecond)
	return ctx, ctxFunc
}

func WatchConfigurationChanges(fileName string) {
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			newHash, _ := utils.CalculateFileHash(fileName)
			if newHash != configurationHash {
				configurationHash = newHash
				LoadConfiguration(fileName)
			}
		}
	}()
}

func Log(dbName, userName, format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, getLogPrefix(dbName, userName)+format+"\n", a...)
}

func LogError(dbName, userName, format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, getLogPrefix(dbName, userName)+"[ERROR] "+format+"\n", a...)
}

func LogAndReturnError(dbName, userName, format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	LogError(dbName, userName, m)
	return errors.New(m)
}

func Trace(dbName, userName, format string, a ...any) {
	if appConfiguration.Trace {
		fmt.Fprintf(gin.DefaultWriter, getLogPrefix(dbName, userName)+"[TRACE] "+format+"\n", a...)
	}
}

func getLogPrefix(dbName, userName string) string {
	return "[" + appConfiguration.AppName + "] [" + dbName + "] [" + userName + "] "
}
