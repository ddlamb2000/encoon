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

var (
	testRouter = gin.Default()
)

func TestSystem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetApiRoutes(testRouter)
	forceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
	t.Run("ConnectDb", func(t *testing.T) { RunTestConnectDbServers(t) })
	t.Run("RecreateDb", func(t *testing.T) { RunTestRecreateDb(t) })
	t.Run("Auth", func(t *testing.T) { RunSystemTestAuth(t) })
	t.Run("Get", func(t *testing.T) { RunSystemTestGet(t) })
	t.Run("Post", func(t *testing.T) { RunSystemTestPost(t) })
	t.Run("DisconnectDb", func(t *testing.T) { RunTestDisconnectDbServers(t) })
}

func getTokenForUser(dbName, userName, userUuid string) string {
	expiration := time.Now().Add(time.Duration(utils.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken(dbName, userName, userUuid, userName, userName, expiration)
	return token
}

func runPOSTRequestForUser(dbName, userName, userUuid, uri, body string) ([]byte, error, int) {
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, err, w.Code
}

func runGETRequestForUser(dbName, userName, userUuid, uri string) ([]byte, error, int) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err, 0
	}
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, err, w.Code
}

func httpCodeEqual(t *testing.T, code int, expectCode int) {
	if code != expectCode {
		t.Errorf(`Response code %v instead of %v.`, code, expectCode)
	}
}

func stringEqual(t *testing.T, got, expect string) {
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func stringNotEqual(t *testing.T, got, expect string) {
	if got == expect {
		t.Errorf(`Got %v.`, got)
	}
}

func byteEqualString(t *testing.T, got []byte, expect string) {
	gotString := string(got)
	if gotString != expect {
		t.Errorf(`Got %v instead of %v.`, gotString, expect)
	}
}

func intEqual(t *testing.T, got, expect int) {
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
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
	utils.LoadConfiguration("../configurations/", "configuration.yml")
	if err := ConnectDbServers(utils.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
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
	if err := DisconnectDbServers(utils.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't disconnect to databases: %v.`, err)
	}
}
