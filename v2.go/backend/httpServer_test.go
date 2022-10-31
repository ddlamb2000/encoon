// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"d.lambert.fr/encoon/utils"
)

func TestHttpServer(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	ConnectDbServers(utils.DatabaseConfigurations)
	setApiRoutes()
	t.Run("404Html", func(t *testing.T) { RunTest404Html(t) })
	t.Run("ApiUsersNoHeader", func(t *testing.T) { RunTestApiUsersNoHeader(t) })
	t.Run("ApiUsersIncorrectToken", func(t *testing.T) { RunTestApiUsersIncorrectToken(t) })
	t.Run("ApiUsersMissingBearer", func(t *testing.T) { RunTestApiUsersMissingBearer(t) })
	t.Run("ApiUsersPassing", func(t *testing.T) { RunTestApiUsersPassing(t) })
	t.Run("ApiUsersNotFound", func(t *testing.T) { RunTestApiUsersNotFound(t) })
	t.Run("ApiUsersNotFound2", func(t *testing.T) { RunTestApiUsersNotFound2(t) })
}

func RunTest404Html(t *testing.T) {
	req, _ := http.NewRequest("GET", "/xxx/yyy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

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

	expected := utils.CleanupStrings(`{ "error": "Not authorized."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersMissingBearer(t *testing.T) {
	token, _ := getNewToken("test", "root", "0", "root", "root")
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

	expected := utils.CleanupStrings(`{ "error": "Unauthorized (invalid request: illegal base64 data at input byte 28)."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersPassing(t *testing.T) {
	ConnectDbServers(utils.DatabaseConfigurations)
	token, _ := getNewToken("test", "root", "0", "root", "root")
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

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
	token, _ := getNewToken("test", "root", "0", "root", "root")
	req, _ := http.NewRequest("GET", "/test/api/v0/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

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
	token, _ := getNewToken("test", "root", "0", "root", "root")
	req, _ := http.NewRequest("GET", "/test/api/v1/us", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

	expected := utils.CleanupStrings(`{ "error": "GRID NOT FOUND"}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}
