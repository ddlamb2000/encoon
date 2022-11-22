// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestGet(t *testing.T) {
	t.Run("VerifyIncorrectUserUuid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", "xxyyzz", "/test/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `{"error":"User not authorized."}`)
	})

	t.Run("VerifyDbNotFound", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/tst/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `{"error":"Invalid request or unauthorized database access: signature is invalid."}`)
	})

	t.Run("VerifyGridNotFound", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `Error when retrieving grid definition`)
	})

	t.Run("VerifyActualRows", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":7`)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `"createdBy":"`+model.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Columns","name":"relationship1","type":"Reference"`)
		jsonStringContains(t, responseData, `"label":"Description","name":"text2","type":"Text"`)
	})

	t.Run("VerifyMissingRow", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids+"/"+model.UuidRootUser)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringDoesntContain(t, responseData, `"countRows"`)
		jsonStringContains(t, responseData, `{"error":"Data not found."}`)
	})

	t.Run("VerifyActualRowSingle", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids+"/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
	})

	t.Run("VerifyActualRowSingleWithTimeOut", func(t *testing.T) {
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		defer database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids+"/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusRequestTimeout)
		jsonStringContains(t, responseData, `{"error":"Get request has been cancelled: context deadline exceeded."}`)
	})

	t.Run("VerifyActualRowSingleBis", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids+"/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
	})
}
