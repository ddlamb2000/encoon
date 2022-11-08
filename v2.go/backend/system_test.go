// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"net/http/httptest"
	"testing"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

func TestSystem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setApiRoutes()
	t.Run("ConnectDb", func(t *testing.T) { RunTestConnectDbServers(t) })
	t.Run("RecreateDb", func(t *testing.T) { RunTestRecreateDb(t) })
	t.Run("Auth", func(t *testing.T) { RunSystemTestAuth(t) })
	t.Run("Post", func(t *testing.T) { RunSystemTestPost(t) })
	t.Run("DisconnectDb", func(t *testing.T) { RunTestDisconnectDbServers(t) })
}

func assertHttpCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	if w.Code != expectedCode {
		t.Errorf(`Response code %v instead of %v.`, w.Code, expectedCode)
	}
}

func RunTestConnectDbServers(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	if err := ConnectDbServers(utils.DatabaseConfigurations); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
	forceTestSleepTime(dbName, 0)
	db := getDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func RunTestRecreateDb(t *testing.T) {
	dbName := "test"
	db := getDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := pingDb(ctx, db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
	if err := recreateDb(ctx, db, dbName); err != nil {
		t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
	}
}

func RunTestDisconnectDbServers(t *testing.T) {
	dbName := "test"
	db := getDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := pingDb(ctx, db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
	if err := DisconnectDbServers(utils.DatabaseConfigurations); err != nil {
		t.Errorf(`Can't disconnect to databases: %v.`, err)
	}
}
