// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"fmt"
	"testing"

	"d.lambert.fr/encoon/backend/utils"
)

func TestSetAndStartHttpServer(t *testing.T) {
	SetAndStartHttpServer()
	if srv.Addr != fmt.Sprintf(":%d", utils.Configuration.HttpServer.Port) {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}
