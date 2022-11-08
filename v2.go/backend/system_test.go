// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func getTokenForUser(dbName, userName, userUuid string) string {
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken(dbName, userName, userUuid, userName, userName, expiration)
	return token
}

func runPOSTRequestForUser(dbName, userName, userUuid, uri, body string) ([]byte, error, int) {
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, err, w.Code
}

func httpCodeEqual(t *testing.T, code int, expectedCode int) {
	if code != expectedCode {
		t.Errorf(`Response code %v instead of %v.`, code, expectedCode)
	}
}

func byteEqualString(t *testing.T, got []byte, expected string) {
	gotString := string(got)
	if gotString != expected {
		t.Errorf(`Got %v instead of %v.`, gotString, expected)
	}
}

func errorIsNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
}

func jsonStringContains(t *testing.T, got []byte, expect string) {
	gotString := utils.CleanupStrings(string(got))
	expectString := utils.CleanupStrings(expect)
	if !strings.Contains(gotString, expectString) {
		t.Errorf(`Response %v doesn't include %v.`, gotString, expectString)
	}
}

func jsonStringDoesntContain(t *testing.T, got []byte, expect string) {
	gotJson := utils.CleanupStrings(string(got))
	expectJson := utils.CleanupStrings(expect)
	if strings.Contains(gotJson, expectJson) {
		t.Errorf(`Response %v includes %v.`, gotJson, expectJson)
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