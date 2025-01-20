// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

var testRouter = gin.Default()

func TestSystem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetApiRoutes(testRouter)
	setDefaultTestSleepTimeAndTimeOutThreshold()
	InitializeCaches()
	t.Run("IncorrectDb", func(t *testing.T) { RunTestConnectDbServersIncorrect(t) })
	t.Run("ConnectDb", func(t *testing.T) { RunTestConnectDbServers(t) })
	t.Run("RecreateDb", func(t *testing.T) { RunTestRecreateDb(t) })
	t.Run("Kafka", func(t *testing.T) { RunSystemTestKafka(t) })
	t.Run("Auth", func(t *testing.T) { RunSystemTestAuth(t) })
	t.Run("Get", func(t *testing.T) { RunSystemTestGet(t) })
	t.Run("Post", func(t *testing.T) { RunSystemTestPost(t) })
	InitializeCaches()
	t.Run("PostRelationships", func(t *testing.T) { RunSystemTestPostRelationships(t) })
	t.Run("RowLevelAccess", func(t *testing.T) { RunSystemTestGetRowLevel(t) })
	t.Run("Cache", func(t *testing.T) { RunSystemTestCache(t) })
	t.Run("NotOwnedColumn", func(t *testing.T) { RunSystemTestNotOwnedColumn(t) })
}

func getTokenForUser(dbName, userName, userUuid string) string {
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken(dbName, userName, userUuid, userName, userName, expiration)
	return token
}

func runPOSTRequestForUser(dbName, userName, userUuid, uri, body string) ([]byte, int, error) {
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, w.Code, err
}

func runGETRequestForUser(dbName, userName, userUuid, uri string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, w.Code, err
}

func httpCodeEqual(t *testing.T, code int, expectCode int) {
	if code != expectCode {
		t.Errorf(`Response code %v instead of %v.`, code, expectCode)
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

func setDefaultTestSleepTimeAndTimeOutThreshold() {
	database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 5000)
}

func RunTestConnectDbServersIncorrect(t *testing.T) {
	configuration.LoadConfiguration("abc/configuration.yml")
	dbName := "alpha"
	_, err := database.GetDbByName(dbName)
	if err == nil {
		t.Errorf(`Database %q found!`, dbName)
	}
}

func RunTestConnectDbServers(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	dbName := "test"
	_, err := database.GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func RunTestRecreateDb(t *testing.T) {
	dbName := "test"
	db, err := database.GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := database.RecreateDb(context.Background(), db, dbName); err != nil {
		t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
	}
}
