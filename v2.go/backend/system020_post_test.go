// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/utils"
)

func RunSystemTestPost(t *testing.T) {
	t.Run("CreateNewUser", func(t *testing.T) { RunTestApiCreateNewUser(t) })
}

func RunTestApiCreateNewUser(t *testing.T) {
	expiration := time.Now().Add(time.Duration(utils.Configuration.HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken("test", "root", utils.UuidRootUser, "root", "root", expiration)

	postStr := `{"rowsAdded":`
	postStr += `[`
	postStr += `{"text01":"aaaa","text02":"bbbb"}`
	postStr += `]`
	postStr += `}`

	req, _ := http.NewRequest("POST", "/test/api/v1/_users", bytes.NewBuffer([]byte(postStr)))
	req.Header.Add("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	assertHttpCode(t, w, http.StatusOK)

	expected := utils.CleanupStrings(`"text01":"aaaa","text02":"bbbb"`)
	response := utils.CleanupStrings(string(responseData))

	if err != nil {
		t.Errorf(`Response %v for %v: %v.`, response, w, err)
	}
	if !strings.Contains(response, expected) {
		t.Errorf(`Response %v incorrect.`, response)
	}
}
