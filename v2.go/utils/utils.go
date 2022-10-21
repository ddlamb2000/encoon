// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Log(message string) {
	fmt.Fprint(gin.DefaultWriter, "[εncooη] ", message, "\n")
}
func Logf(format string, v ...any) {
	fmt.Fprintf(gin.DefaultWriter, "[εncooη] "+format, v)
}

func LogFatal(v ...any) {
	fmt.Fprint(gin.DefaultWriter, "[εncooη] ", v, "\n")
	log.Fatal("[εncooη]", v)
}
