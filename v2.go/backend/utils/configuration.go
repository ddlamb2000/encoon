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

func LoadConfiguration(directory string) {
	loadMainConfiguration(directory, "configuration.yml")
	loadDatabaseConfigurations(directory, "databases/")
}

func loadMainConfiguration(directory string, fileName string) {
	f, err := os.Open(directory + fileName)
	if err != nil {
		LogError("Error loading configuration from file %q: %v", fileName, err)
		return
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Configuration)
	if err != nil {
		LogError("Error parsing configuration from file %q:", fileName, err)
	} else {
		Log("Configuration loaded from file %q:", fileName)
	}
}

func loadDatabaseConfigurations(directory string, subDirectory string) {
	files, err := ioutil.ReadDir(directory + subDirectory)
	if err != nil {
		LogError("Load configuration: %v", err)
		return
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "yml") {
			loadDatabaseConfiguration(directory + subDirectory + file.Name())
		}
	}
}

func loadDatabaseConfiguration(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		LogError("Error loading configuration from file %q: %v", fileName, err)
		return
	}
	defer f.Close()
	var databaseConfiguration DatabaseConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&databaseConfiguration)
	if err != nil {
		LogError("Error parsing configuration from file %q:", fileName, err)
	} else {
		dbName := databaseConfiguration.Database.Name
		Log("Load database %q configuration from file %q:", dbName, fileName)
		DatabaseConfigurations[dbName] = &databaseConfiguration
	}
}

func IsDatabaseEnabled(dbName string) bool {
	return dbName != "" && DatabaseConfigurations[dbName] != nil
}

func GetJWTSecret(dbName string) []byte {
	return []byte(dbName + DatabaseConfigurations[dbName].Database.JwtSecret)
}

func GetRootAndPassowrd(dbName string) (string, string) {
	return DatabaseConfigurations[dbName].Database.Root, DatabaseConfigurations[dbName].Database.Password
}
