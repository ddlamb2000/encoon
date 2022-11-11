// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package configuration

import (
	"context"
	"os"
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
}

var appConfiguration Configuration

func LoadConfiguration(directory, fileName string) error {
	f, err := os.Open(directory + fileName)
	if err != nil {
		utils.LogError("Error loading configuration from file %q: %v.", fileName, err)
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&appConfiguration)
	if err != nil {
		utils.LogError("Error parsing configuration from file %q: %v.", fileName, err)
		return err
	}
	utils.Log("Configuration loaded from file %q.", fileName)
	return nil
}

func GetConfiguration() Configuration { return appConfiguration }

func loadMainConfiguration(directory string, fileName string) error {
	f, err := os.Open(directory + fileName)
	if err != nil {
		utils.LogError("Error loading configuration from file %q: %v.", fileName, err)
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&appConfiguration)
	if err != nil {
		utils.LogError("Error parsing configuration from file %q: %v.", fileName, err)
		return err
	}
	utils.Log("Configuration loaded from file %q.", fileName)
	return nil
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
	dbConfiguration := GetDatabaseConfiguration(dbName)
	threshold := dbConfiguration.TimeOutThreshold
	if threshold < 10 {
		threshold = 10
	}
	ctx, ctxFunc := context.WithTimeout(context.Background(), time.Duration(threshold)*time.Millisecond)
	return ctx, ctxFunc
}
