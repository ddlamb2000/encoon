// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpServer struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"httpServer"`
	Database struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"database"`
}

var (
	Configuration Config
)

func Log(message string) {
	fmt.Fprint(gin.DefaultWriter, "[εncooη] "+message+"\n")
}

func LogFatal(v ...any) {
	fmt.Fprint(gin.DefaultWriter, "[εncooη] ", v, "\n")
	log.Fatal("[εncooη]", v)
}

func InitWithLog() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	Log("Starting.")

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
}
