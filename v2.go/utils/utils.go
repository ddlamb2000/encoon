// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Log(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, format+"\n", a...)
}

func Trace(trace, format string, a ...any) {
	if trace != "" {
		fmt.Fprintf(gin.DefaultWriter, "[TRACE] "+format+"\n", a...)
	}
}

func LogError(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "[ERROR] "+format+"\n", a...)
}

func LogAndReturnError(format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	LogError(m)
	return errors.New(m)
}

func CleanupStrings(s string) string {
	return strings.Join(strings.Fields(strings.Replace(s, "\n", "", -1)), " ")
}
