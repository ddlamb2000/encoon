// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
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
		defer database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		req, _ := http.NewRequest("POST", "/test/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "dGVzdA=="}`))
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)

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
		token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
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
		token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
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
		token, _ := getNewToken("test", "root", model.UuidRootUser, "root", "root", expiration)
		req, _ := http.NewRequest("GET", "/test/api/v1/_users", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)
		responseData, err := io.ReadAll(w.Body)
		httpCodeEqual(t, w.Code, http.StatusOK)

		expect := utils.CleanupStrings(`"text1":"root","text2":"root"`)
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
		token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
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
		token, _ := getNewToken("test", "root", model.UuidRootUser, "root", "root", expiration)
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

	t.Run("CreateNewUserNoAuth", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"aaaa","text2":"bbbb"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("te st", "root", model.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringContains(t, responseData, `"error":"Invalid request or unauthorized database access: signature is invalid."`)
	})

	t.Run("Post404", func(t *testing.T) {
		postStr := `{}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		byteEqualString(t, responseData, `404 page not found`)
	})

	t.Run("CreateUserNoData", func(t *testing.T) {
		postStr := `{}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+model.UuidUsers+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `{"uuid":"`+model.UuidRootUser+`"`)
	})

	t.Run("CreateNewSingleUser", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test01","text2":"Zero-one","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringDoesntContain(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"countRows":2`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+model.UuidUsers+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `{"uuid":"`+model.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"Zero-one","text3":"Test"`)
	})

	t.Run("Create3Users", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test02","text2":"Zero-two","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"},` +
			`{"text1":"test03","text2":"Zero-three","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"},` +
			`{"text1":"test04","text2":"Zero-four","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringDoesntContain(t, responseData, `"countRows":2`)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringContains(t, responseData, `"text1":"test02","text2":"Zero-two","text3":"Test"`)
		jsonStringContains(t, responseData, `"text1":"test03","text2":"Zero-three","text3":"Test"`)
		jsonStringContains(t, responseData, `"text1":"test04","text2":"Zero-four","text3":"Test"`)
	})

	t.Run("CreateWithIncorrectUserUuid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test02","text2":"Zero-two","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "root", "xxyyzz", "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
	})
}
