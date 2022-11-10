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
	fmt.Fprintf(gin.DefaultWriter, "["+Configuration.AppName+"] "+format+"\n", a...)
}

func LogError(format string, a ...any) {
	fmt.Fprintf(gin.DefaultWriter, "["+Configuration.AppName+"] [ERROR] "+format+"\n", a...)
}

func LogAndReturnError(format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	LogError(m)
	return errors.New(m)
}

func CleanupStrings(s string) string {
	return strings.Join(strings.Fields(strings.Replace(s, "\n", "", -1)), " ")
}

func GetContextWithTimeOut() (context.Context, context.CancelFunc) {
	threshold := Configuration.TimeOutThreshold
	if threshold == 0 {
		threshold = 10
	}
	ctx, ctxFunc := context.WithTimeout(context.Background(), time.Duration(threshold)*time.Millisecond)
	return ctx, ctxFunc
}
