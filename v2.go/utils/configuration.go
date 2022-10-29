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
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"httpServer"`
}

type DatabaseConfig struct {
	Database struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Name      string `yaml:"name"`
		User      string `yaml:"user"`
		JwtSecret string `yaml:"jwtsecret"`
		Root      string `yaml:"root"`
		Password  string `yaml:"password"`
	} `yaml:"database"`
}

var (
	Configuration          Config
	DatabaseConfigurations = make(map[string]*DatabaseConfig)
)

func LoadConfiguration(directory string) bool {
	if loadMainConfiguration(directory, "configuration.yml") {
		return true
	}
	if loadDatabaseConfigurations(directory, "databases/") {
		return true
	}
	return false
}

func loadMainConfiguration(directory string, fileName string) bool {
	f, err := os.Open(directory + fileName)
	if err != nil {
		LogError("Error loading configuration from file %q: %v.", fileName, err)
		return true
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Configuration)
	if err != nil {
		LogError("Error parsing configuration from file %q: %v.", fileName, err)
		return true
	}
	Log("Configuration loaded from file %q.", fileName)
	return false
}

func loadDatabaseConfigurations(directory string, subDirectory string) bool {
	files, err := ioutil.ReadDir(directory + subDirectory)
	if err != nil {
		LogError("Load configuration: %v.", err)
		return true
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "yml") {
			if loadDatabaseConfiguration(directory + subDirectory + file.Name()) {
				return true
			}
		}
	}
	return false
}

func loadDatabaseConfiguration(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil {
		LogError("Error loading configuration from file %q: %v", fileName, err)
		return true
	}
	defer f.Close()
	var databaseConfiguration DatabaseConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&databaseConfiguration)
	if err != nil {
		LogError("Error parsing configuration from file %q:", fileName, err)
		return true
	}
	dbName := databaseConfiguration.Database.Name
	Log("Load database %q configuration from file %q:", dbName, fileName)
	DatabaseConfigurations[dbName] = &databaseConfiguration
	return false
}

func IsDatabaseEnabled(dbName string) bool {
	return dbName != "" && DatabaseConfigurations[dbName] != nil
}

func GetJWTSecret(dbName string) []byte {
	return []byte(dbName + DatabaseConfigurations[dbName].Database.JwtSecret)
}

func GetRootAndPassword(dbName string) (string, string) {
	return DatabaseConfigurations[dbName].Database.Root, DatabaseConfigurations[dbName].Database.Password
}
