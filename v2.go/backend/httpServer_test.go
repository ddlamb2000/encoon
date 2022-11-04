// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
)

func TestHttpServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.LoadConfiguration("../configurations/")
	ConnectDbServers(utils.DatabaseConfigurations)
	setApiRoutes()
	t.Run("AuthInvalid1", func(t *testing.T) { RunTestAuthInvalid1(t) })
	t.Run("AuthInvalid2", func(t *testing.T) { RunTestAuthInvalid2(t) })
	t.Run("AuthValid", func(t *testing.T) { RunTestAuthValid(t) })
	t.Run("404Html", func(t *testing.T) { RunTest404Html(t) })
	t.Run("ApiUsersNoHeader", func(t *testing.T) { RunTestApiUsersNoHeader(t) })
	t.Run("ApiUsersIncorrectToken", func(t *testing.T) { RunTestApiUsersIncorrectToken(t) })
	t.Run("ApiUsersIncorrectToken2", func(t *testing.T) { RunTestApiUsersIncorrectToken2(t) })
	t.Run("ApiUsersMissingBearer", func(t *testing.T) { RunTestApiUsersMissingBearer(t) })
	t.Run("ApiUsersExpired", func(t *testing.T) { RunTestApiUsersExpired(t) })
	t.Run("ApiUsersPassing", func(t *testing.T) { RunTestApiUsersPassing(t) })
	t.Run("ApiUsersNotFound", func(t *testing.T) { RunTestApiUsersNotFound(t) })
	t.Run("ApiUsersNotFound2", func(t *testing.T) { RunTestApiUsersNotFound2(t) })
}

func assertHttpCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	if w.Code != expectedCode {
		t.Errorf(`Response code %v instead of %v.`, w.Code, expectedCode)
	}
}

func RunTestAuthInvalid1(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test/api/v1/authentication", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusBadRequest)

	expected := utils.CleanupStrings(`{"error":"[test] Invalid ID or password: sql: no rows in result set"}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect instead of %v.`, response, expected)
	}
}

func RunTestAuthInvalid2(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "======"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusBadRequest)

	expected := utils.CleanupStrings(`{"error":"[test] Invalid ID or password: sql: no rows in result set"}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect instead of %v.`, response, expected)
	}
}

func RunTestAuthValid(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "dGVzdA=="}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusOK)

	expected := utils.CleanupStrings(`{"token":`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTest404Html(t *testing.T) {
	req, _ := http.NewRequest("GET", "/xxx/yyy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusNotFound)

	expected := utils.CleanupStrings(`404 page not found`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersNoHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusUnauthorized)

	expected := utils.CleanupStrings(`{ "error": "No authorization found."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersIncorrectToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "xxxxxxxxxxx")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusUnauthorized)

	expected := utils.CleanupStrings(`Not authorized for /test/api/v1/users`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect instead of %v.`, response, expected)
	}
}

func RunTestApiUsersIncorrectToken2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "xxxxxxx")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusUnauthorized)

	expected := utils.CleanupStrings(`{ "error": "Incorrect header."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersMissingBearer(t *testing.T) {
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusUnauthorized)

	expected := utils.CleanupStrings(`Invalid request`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersExpired(t *testing.T) {
	ConnectDbServers(utils.DatabaseConfigurations)
	expiration := time.Now().Add(-time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusUnauthorized)

	expected := utils.CleanupStrings(`{ "error": "Authorization expired.", "expired": true}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersPassing(t *testing.T) {
	ConnectDbServers(utils.DatabaseConfigurations)
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusOK)

	expected := utils.CleanupStrings(`"text01": "root", "text02": "root"`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersNotFound(t *testing.T) {
	ConnectDbServers(utils.DatabaseConfigurations)
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
	req, _ := http.NewRequest("GET", "/test/api/v0/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusNotFound)

	expected := utils.CleanupStrings(`404 page not found`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersNotFound2(t *testing.T) {
	ConnectDbServers(utils.DatabaseConfigurations)
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
	req, _ := http.NewRequest("GET", "/test/api/v1/us", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusNotFound)

	expected := utils.CleanupStrings(`{ "error": "[test] [root] Grid us not found."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect instead of %v.`, response, expected)
	}
}
