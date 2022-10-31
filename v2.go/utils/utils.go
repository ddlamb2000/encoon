// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Log(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "["+Configuration.AppName+"] "+format+"\n", a...)
}

func LogError(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "["+Configuration.AppName+"] [ERROR] "+format+"\n", a...)
}

func InitWithLog() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	Log("Starting.")
}

func CleanupStrings(s string) string {
	return strings.Join(strings.Fields(strings.Replace(s, "\n", "", -1)), " ")
}
