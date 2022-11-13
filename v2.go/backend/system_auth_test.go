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

	"d.lambert.fr/encoon/authentication"
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
)

func RunSystemTestAuth(t *testing.T) {
	t.Run("AuthInvalid1", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/test/api/v1/authentication", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`{"error":"[test] Invalid username or passphrase for \"\": sql: no rows in result set."}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if response != expect {
			t.Errorf(`Response %v incorrect instead of %v.`, response, expect)
		}
	})

	t.Run("AuthInvalid2", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/test/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "======"}`))
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`{"error":"[test] Invalid username or passphrase for \"root\": sql: no rows in result set."}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if response != expect {
			t.Errorf(`Response %v incorrect instead of %v.`, response, expect)
		}
	})

	t.Run("AuthValid", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/test/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "dGVzdA=="}`))
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusOK)

		expect := utils.CleanupStrings(`{"token":`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersWithTimeOut", func(t *testing.T) {
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		req, _ := http.NewRequest("POST", "/test/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "dGVzdA=="}`))
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)

		httpCodeEqual(t, w.Code, http.StatusRequestTimeout)

		expect := utils.CleanupStrings(`{"error":"[test] [root] Authentication request has been cancelled: context deadline exceeded."}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("404Html", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/xxx/yyy", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusNotFound)

		expect := utils.CleanupStrings(`404 page not found`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if response != expect {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersNoHeader", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`{"error":"No authorization found."}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if response != expect {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersIncorrectToken", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		req.Header.Add("Authorization", "xxxxxxxxxxx")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`Not authorized for /test/api/v1/_users`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect instead of %v.`, response, expect)
		}
	})

	t.Run("ApiUsersIncorrectToken2", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		req.Header.Add("Authorization", "xxxxxxx")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`{"error":"Incorrect header."}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if response != expect {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersMissingBearer", func(t *testing.T) {
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		token, _ := authentication.GetNewToken("test", "root", "0", "root", "root", expiration)
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		req.Header.Add("Authorization", token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`Invalid request`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersExpired", func(t *testing.T) {
		expiration := time.Now().Add(-time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		token, _ := authentication.GetNewToken("test", "root", "0", "root", "root", expiration)
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		expect := utils.CleanupStrings(`{"error":"Authorization expired.","expired":true}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersPassing", func(t *testing.T) {
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		token, _ := authentication.GetNewToken("test", "root", utils.UuidRootUser, "root", "root", expiration)
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusOK)

		expect := utils.CleanupStrings(`"text01":"root","text02":"root"`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersNotFound", func(t *testing.T) {
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		token, _ := authentication.GetNewToken("test", "root", "0", "root", "root", expiration)
		req, _ := http.NewRequest("GET", "/test/api/v0/users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusNotFound)

		expect := utils.CleanupStrings(`404 page not found`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect.`, response)
		}
	})

	t.Run("ApiUsersNotFound2", func(t *testing.T) {
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		token, _ := authentication.GetNewToken("test", "root", utils.UuidRootUser, "root", "root", expiration)
		req, _ := http.NewRequest("GET", "/test/api/v1/us", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusNotFound)

		expect := utils.CleanupStrings(`{"error":"[test] [root] Grid \"us\" not found."}`)
		response := utils.CleanupStrings(string(responseData))

		if err != nil {
			t.Errorf(`Response %v for %v: %v.`, response, w, err)
		}
		if !strings.Contains(response, expect) {
			t.Errorf(`Response %v incorrect instead of %v.`, response, expect)
		}
	})
}
