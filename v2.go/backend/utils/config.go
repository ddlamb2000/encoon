// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"os"
	"regexp"

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

func LoadConfiguration() {
	f, err := os.Open("configuration.yml")
	if err != nil {
		LogError("Load configuration:", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Configuration)
	if err != nil {
		LogError("Error parsing configuration:", err)
	}

	Configuration.DatabaseNames = regexp.MustCompile(`\s*,\s*`).Split(Configuration.Database.Names, -1)
}

func IsDatabaseEnabled(str string) bool {
	return isNotEmptyAndHasAnyInCommon(Configuration.DatabaseNames, str)
}
