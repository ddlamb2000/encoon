// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpServer struct {
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		JwtExpiration int    `yaml:"jwtExpiration"`
	} `yaml:"httpServer"`

	AppName   string `yaml:"appName"`
	AppTag    string `yaml:"appTag"`
	DbTimeOut int    `yaml:"dbTimeOut"`
}

type DatabaseConfig struct {
	Database struct {
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		Name          string `yaml:"name"`
		User          string `yaml:"user"`
		JwtSecret     string `yaml:"jwtsecret"`
		Root          string `yaml:"root"`
		Password      string `yaml:"password"`
		TestSleepTime int    `yaml:"testSleepTime"`
	} `yaml:"database"`
}

var (
	Configuration          Config
	DatabaseConfigurations = make(map[string]*DatabaseConfig)
)

func LoadConfiguration(directory string) error {
	if err := loadMainConfiguration(directory, "configuration.yml"); err != nil {
		return err
	}
	if err := loadDatabaseConfigurations(directory, "databases/"); err != nil {
		return err
	}
	return nil
}

func loadMainConfiguration(directory string, fileName string) error {
	f, err := os.Open(directory + fileName)
	if err != nil {
		LogError("Error loading configuration from file %q: %v.", fileName, err)
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Configuration)
	if err != nil {
		LogError("Error parsing configuration from file %q: %v.", fileName, err)
		return err
	}
	Log("Configuration loaded from file %q.", fileName)
	return nil
}

func loadDatabaseConfigurations(directory string, subDirectory string) error {
	files, err := ioutil.ReadDir(directory + subDirectory)
	if err != nil {
		LogError("Load configuration: %v.", err)
		return err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "yml") {
			if err := loadDatabaseConfiguration(directory + subDirectory + file.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadDatabaseConfiguration(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		LogError("Error loading configuration from file %q: %v", fileName, err)
		return err
	}
	defer f.Close()
	var databaseConfiguration DatabaseConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&databaseConfiguration)
	if err != nil {
		LogError("Error parsing configuration from file %q:", fileName, err)
		return err
	}
	dbName := databaseConfiguration.Database.Name
	Log("Load database %q configuration from file %q:", dbName, fileName)
	DatabaseConfigurations[dbName] = &databaseConfiguration
	return nil
}

func IsDatabaseEnabled(dbName string) bool {
	return dbName != "" && DatabaseConfigurations[dbName] != nil
}

func GetJWTSecret(dbName string) []byte {
	if !IsDatabaseEnabled(dbName) {
		return nil
	}
	return []byte(dbName + DatabaseConfigurations[dbName].Database.JwtSecret)
}

func GetRootAndPassword(dbName string) (string, string) {
	if !IsDatabaseEnabled(dbName) {
		return "", ""
	}
	return DatabaseConfigurations[dbName].Database.Root, DatabaseConfigurations[dbName].Database.Password
}
