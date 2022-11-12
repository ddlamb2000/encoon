// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Log(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, format+"\n", a...)
}

func LogError(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "[ERROR] "+format+"\n", a...)
}

func LogAndReturnError(format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	LogError(m)
	return errors.New(m)
}

func Trace(trace, format string, a ...any) {
	if trace != "" {
		fmt.Fprintf(gin.DefaultWriter, "[TRACE] "+format+"\n", a...)
	}
}

func CleanupStrings(s string) string {
	return strings.Join(strings.Fields(strings.Replace(s, "\n", "", -1)), " ")
}

func CalculateFileHash(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		LogError("Can't open file %v: %v", fileName, err)
		return ""
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		LogError("Can't read file %v: %v", fileName, err)
		return ""
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum
}
