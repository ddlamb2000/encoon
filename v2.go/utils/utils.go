// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Log(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "["+appConfiguration.AppName+"] "+format+"\n", a...)
}

func Trace(trace, format string, a ...any) {
	if trace != "" {
		fmt.Fprintf(gin.DefaultWriter, "["+appConfiguration.AppName+"] [TRACE] "+format+"\n", a...)
	}
}

func LogError(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "["+appConfiguration.AppName+"] [ERROR] "+format+"\n", a...)
}

func LogAndReturnError(format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	LogError(m)
	return errors.New(m)
}

func CleanupStrings(s string) string {
	return strings.Join(strings.Fields(strings.Replace(s, "\n", "", -1)), " ")
}

func GetContextWithTimeOut(dbName string) (context.Context, context.CancelFunc) {
	dbConfiguration := GetDatabaseConfiguration(dbName)
	threshold := dbConfiguration.TimeOutThreshold
	if threshold < 10 {
		threshold = 10
	}
	ctx, ctxFunc := context.WithTimeout(context.Background(), time.Duration(threshold)*time.Millisecond)
	return ctx, ctxFunc
}
