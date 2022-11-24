// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"errors"
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
		jsonStringContains(t, responseData, `"error":"Data not found."`)
	})

	t.Run("VerifyActualRows", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `"createdBy":"`+model.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Columns","name":"relationship1","type":"Reference"`)
		jsonStringContains(t, responseData, `"label":"Description","name":"text2","type":"Text"`)
	})

	t.Run("VerifyDbNotConfigured", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("baddb", "root", model.UuidRootUser, "/baddb/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringContains(t, responseData, `Unable to connect to database: dial tcp`)
	})

	t.Run("VerifyDbNotConfigured2", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/undefined/api/v1/xxx")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringContains(t, responseData, `{"error":"No database parameter."}`)
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

	t.Run("VerifyActualRowsWithDefect1", func(t *testing.T) {
		getRowsQueryForGridsApiImpl := getRowsQueryForGridsApi
		getRowsQueryForGridsApi = func(grid *model.Grid, uuid string) string { return "xxx" } // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when querying rows: pq: syntax error`)
		getRowsQueryForGridsApi = getRowsQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect2", func(t *testing.T) {
		getGridQueryForGridsApiImpl := getGridQueryForGridsApi
		getGridQueryForGridsApi = func() string { return "xxx" } // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when retrieving grid definition: pq: syntax error`)
		getGridQueryForGridsApi = getGridQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect3", func(t *testing.T) {
		getGridColumsQueryForGridsApiImpl := getGridColumsQueryForGridsApi
		getGridColumsQueryForGridsApi = func() string { return "xxx" } // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when querying columns: pq: syntax error`)
		getGridColumsQueryForGridsApi = getGridColumsQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect4", func(t *testing.T) {
		getGridColumnQueryOutputForGridsApiImpl := getGridColumnQueryOutputForGridsApi
		getGridColumnQueryOutputForGridsApi = func(column *model.Column) []any { return nil } // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when scanning columns for: sql`)
		getGridColumnQueryOutputForGridsApi = getGridColumnQueryOutputForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect5", func(t *testing.T) {
		getQueryReferencedRowsForRowImpl := getQueryReferencedRowsForRow
		getQueryReferencedRowsForRow = func() string { return "xxx" } // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when querying referenced rows: pq: syntax error`)
		getQueryReferencedRowsForRow = getQueryReferencedRowsForRowImpl
	})

	t.Run("VerifyActualRowsWithDefect6", func(t *testing.T) {
		getQueryReferencedRowsForRowImpl := getQueryReferencedRowsForRow
		getQueryReferencedRowsForRow = func() string {
			return "SELECT NULL, NULL FROM rows WHERE gridUuid = $1 AND text1 = $2 AND text2 = $3 AND text3 = $4"
		} // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when scanning referenced rows: sql`)
		getQueryReferencedRowsForRow = getQueryReferencedRowsForRowImpl
	})

	t.Run("VerifyActualRowsWithDefect7", func(t *testing.T) {
		getGridForGridsApiImpl := getGridForGridsApi
		getGridForGridsApi = func(r apiRequestParameters, gridUuid string) (*model.Grid, error) { return nil, errors.New("xxx") } // mock function
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `"error":"xxx"`)
		getGridForGridsApi = getGridForGridsApiImpl
	})

	t.Run("VerifyActualRowSingleDefect", func(t *testing.T) {
		getRowsQueryParametersForGridsApiImpl := getRowsQueryParametersForGridsApi
		getRowsQueryParametersForGridsApi = func(gridUuid, uuid string) []any { return nil }
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids+"/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when querying rows: pq`)
		getRowsQueryParametersForGridsApi = getRowsQueryParametersForGridsApiImpl
	})

	t.Run("VerifyActualRowSingleDefect", func(t *testing.T) {
		getRowsQueryOutputForGridsApiImpl := getRowsQueryOutputForGridsApi
		getRowsQueryOutputForGridsApi = func(grid *model.Grid, row *model.Row) []any { return nil }
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids+"/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Error when scanning rows: sql
		`)
		getRowsQueryOutputForGridsApi = getRowsQueryOutputForGridsApiImpl
	})
}
