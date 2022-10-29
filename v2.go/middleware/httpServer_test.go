// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"d.lambert.fr/encoon/backend"
	"d.lambert.fr/encoon/utils"
)

func TestApiUsers(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	backend.LoadData()
	setApiRoutes()
	t.Run("TestApiUsersNoHeader", func(t *testing.T) { RunTestApiUsersNoHeader(t) })
	t.Run("TestApiUsersIncorrectToken", func(t *testing.T) { RunTestApiUsersIncorrectToken(t) })
	t.Run("TestApiUsersMissingBearer", func(t *testing.T) { RunTestApiUsersMissingBearer(t) })
	t.Run("TestApiUsersPassing", func(t *testing.T) { RunTestApiUsersPassing(t) })
}

func RunTestApiUsersNoHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

	expected := utils.CleanupStrings(`{ "message": "Not authorized."}`)
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

	expected := utils.CleanupStrings(`{ "message": "Not authorized."}`)
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

	expected := utils.CleanupStrings(`{ "disconnect": true, "error": true, "message": "Unauthorized."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Errorf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersPassing(t *testing.T) {
	token, _ := getNewToken("test", "root", "0", "root", "root")
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

	expected := utils.CleanupStrings(`"version": 1, "enabled": true, "email": "root@encoon.com", "firstName": "Root", "lastName": "Encoon"`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}
