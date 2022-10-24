// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func Log(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "[εncooη] "+format+"\n", a...)
}

func LogError(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "[εncooη][ERROR] "+format+"\n", a...)
}

func InitWithLog() {
	f, _ := os.Create("logs/encoon.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	Log("Starting.")
}

func isNotEmptyAndHasAnyInCommon(strs []string, str string) bool {
	if str == "" {
		return false
	}
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}
