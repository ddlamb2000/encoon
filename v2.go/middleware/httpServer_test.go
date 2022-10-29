// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"d.lambert.fr/encoon/backend"
	"d.lambert.fr/encoon/utils"
)

func TestApiUsers(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	backend.LoadData()
	setApiRoutes()
	RunTestApiUsersNoHeader(t)
	RunTestApiUsersIncorrectToken(t)
	RunTestApiUsersMissingBearer(t)
	RunTestApiUsersPassing(t)
}

func RunTestApiUsersNoHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

	expected := utils.CleanupStrings(`{ "message": "Not authorized."}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Fatalf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Fatalf(`Response %v incorrect.`, response)
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
		t.Fatalf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Fatalf(`Response %v incorrect.`, response)
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
		t.Fatalf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Fatalf(`Response %v incorrect.`, response)
	}
}

func RunTestApiUsersPassing(t *testing.T) {
	token, _ := getNewToken("test", "root", "0", "root", "root")
	req, _ := http.NewRequest("GET", "/test/api/v1/users", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)

	expected := utils.CleanupStrings(`{ "users": [ { "uuid": "c788a76d-4aa6-4073-8904-35a9b99a3289", "version": 1, "enabled": true, "email": "root@encoon.com", "firstName": "Root", "lastName": "Encoon" }, { "uuid": "bced42a2-6ddd-4023-ad40-0d46962b7872", "version": 1, "enabled": true, "email": "system@encoon.com", "firstName": "System", "lastName": "Encoon" }, { "uuid": "67b560b9-63ff-4fed-9b64-26c7f86e540c", "version": 0, "enabled": false, "email": "none@encoon.com", "firstName": "", "lastName": "" } ]}`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Fatalf(`Response %v for %v: %v.`, response, w, err)
	}
	if response != expected {
		t.Fatalf(`Response %v incorrect.`, response)
	}
}
