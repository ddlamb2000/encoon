// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
}
