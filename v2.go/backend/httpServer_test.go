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

func TestApis(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	setApiRoutes()
	t.Run("TestApiUsersNoHeader", func(t *testing.T) { TestApiUsersNoHeader(t) })
	t.Run("TestApiUsersIncorrectToken", func(t *testing.T) { TestApiUsersIncorrectToken(t) })
	t.Run("TestApiUsersMissingBearer", func(t *testing.T) { TestApiUsersMissingBearer(t) })
	t.Run("TestApiUsersPassing", func(t *testing.T) { TestApiUsersPassing(t) })
}

func TestApiUsersNoHeader(t *testing.T) {
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

func TestApiUsersIncorrectToken(t *testing.T) {
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

func TestApiUsersMissingBearer(t *testing.T) {
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

func TestApiUsersPassing(t *testing.T) {
	if err := ConnectDbServers(utils.DatabaseConfigurations); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
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
