// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
)

func RunSystemTestGet(t *testing.T) {
	t.Run("VerifyIncorrectUserUuid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", "xxyyzz", "/test/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `{"error":"User not authorized for /test/api/v1/xxx."}`)
	})

	t.Run("VerifyDbNotFound", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/tst/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `{"error":"Invalid request or unauthorized database access: signature is invalid."}`)
	})

	t.Run("VerifyGridNotFound", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `{"error":"[test] [root] Grid \"xxx\" not found."}`)
	})

	t.Run("VerifyActualRows", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+utils.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `"createdBy":"`+utils.UuidRootUser+`"`)
	})

	t.Run("VerifyMissingRow", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids/"+utils.UuidRootUser)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringDoesntContain(t, responseData, `"countRows"`)
		jsonStringContains(t, responseData, `{"error":"[test] [root] Data not found."}`)
	})

	t.Run("VerifyActualRowSingle", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids/"+utils.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+utils.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[{"uuid":"`+utils.UuidGrids+`"`)
	})

	t.Run("VerifyActualRowSingleWithTimeOut", func(t *testing.T) {
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		defer database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids/"+utils.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusRequestTimeout)
		jsonStringContains(t, responseData, `{"error":"[test] [root] Get request has been cancelled: context deadline exceeded."}`)
	})

	t.Run("VerifyActualRowSingleBis", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids/"+utils.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+utils.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[{"uuid":"`+utils.UuidGrids+`"`)
	})
}
