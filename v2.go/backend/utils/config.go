// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpServer struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"httpServer"`

	Database struct {
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Names string `yaml:"names"`
		User  string `yaml:"user"`
	} `yaml:"database"`

	DatabaseNames []string
}

var (
	Configuration Config
)

func loadConfiguration() {
	f, err := os.Open("configuration.yml")
	if err != nil {
		LogFatal("Load configuration:", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Configuration)
	if err != nil {
		LogFatal("Error parsing configuration:", err)
	}

	Configuration.DatabaseNames = strings.Split(Configuration.Database.Names, ",")
	for i, v := range Configuration.DatabaseNames {
		Configuration.DatabaseNames[i] = strings.Trim(v, " ")
	}
}

func DatabaseAllowed(str string) bool {
	if str == "" {
		return false
	}
	for _, v := range Configuration.DatabaseNames {
		if v == str {
			return true
		}
	}
	return false
}
