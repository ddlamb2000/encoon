// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"errors"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestGet(t *testing.T) {
	t.Run("VerifyGridNotFound", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "d7c004ff-cccc-dddd-eeee-cd42b2847508", ApiParameters{
			Action:   ActionLoad,
			GridUuid: "d7c004ff-cccc-dddd-eeee-cd42b2847508",
		})
		responseIsFailure(t, response)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `"textMessage":"Data not found"`)
	})

	t.Run("VerifyGridNotFound2", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "xxx", ApiParameters{
			Action:   ActionLoad,
			GridUuid: "xxx",
		})
		responseIsFailure(t, response)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows":`)
		jsonStringContains(t, responseData, `Error when retrieving grid definition: pq: invalid input syntax for type uuid`)
	})

	t.Run("VerifyActualRows", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `"createdBy":"`+model.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Columns","name":"relationship1","type":"Reference"`)
		jsonStringContains(t, responseData, `"label":"Description","name":"text2","type":"Text"`)
	})

	t.Run("VerifyDbNotConfigured", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "baddb", "root", model.UuidRootUser, "xxx", ApiParameters{
			Action:   ActionLoad,
			GridUuid: "xxx",
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Unable to connect to database: dial tcp`)
	})

	t.Run("VerifyMissingRow", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidRootUser,
		})
		responseIsFailure(t, response)
		jsonStringDoesntContain(t, responseData, `"countRows"`)
		jsonStringContains(t, responseData, `"textMessage":"Data not found"`)
	})

	t.Run("VerifyActualRowSingle", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
	})

	t.Run("VerifyActualRowSingleWithTimeOut", func(t *testing.T) {
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		defer setDefaultTestSleepTimeAndTimeOutThreshold()
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Get request has been cancelled: context deadline exceeded."`)
	})

	t.Run("VerifyActualRowSingleBis", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
	})

	t.Run("VerifyActualRowsWithDefect1", func(t *testing.T) {
		getRowsQueryForGridsApiImpl := getRowsQueryForGridsApi
		getRowsQueryForGridsApi = func(*model.Grid, string, bool) string { return "xxx" } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when querying rows: pq: syntax error`)
		getRowsQueryForGridsApi = getRowsQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect2", func(t *testing.T) {
		getGridQueryForGridsApiImpl := getGridQueryForGridsApi
		getGridQueryForGridsApi = func() string { return "xxx" } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when retrieving grid definition: pq: syntax error`)
		getGridQueryForGridsApi = getGridQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect3", func(t *testing.T) {
		getGridColumsOwnedQueryForGridsApiImpl := getGridColumsOwnedQueryForGridsApi
		getGridColumsOwnedQueryForGridsApi = func() string { return "xxx" } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when querying columns: pq: syntax error`)
		getGridColumsOwnedQueryForGridsApi = getGridColumsOwnedQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect4", func(t *testing.T) {
		getGridColumnQueryOutputForGridsApiImpl := getGridColumnQueryOutputForGridsApi
		getGridColumnQueryOutputForGridsApi = func(column *model.Column) []any { return nil } // mock function
		response, _ := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		getGridColumnQueryOutputForGridsApi = getGridColumnQueryOutputForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect5", func(t *testing.T) {
		getQueryReferencedRowsForRowImpl := getQueryReferencedRowsForRow
		getQueryReferencedRowsForRow = func(bool) string { return "xxx" } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when querying referenced rows: pq: syntax error`)
		getQueryReferencedRowsForRow = getQueryReferencedRowsForRowImpl
	})

	t.Run("VerifyActualRowsWithDefect6", func(t *testing.T) {
		getQueryReferencedRowsForRowImpl := getQueryReferencedRowsForRow
		getQueryReferencedRowsForRow = func(bool) string {
			return "SELECT NULL, NULL FROM relationships WHERE gridUuid = $1 AND text1 = $2 AND text2 = $3 AND text3 = $4 AND text2 = $5"
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when scanning referenced rows: sql`)
		getQueryReferencedRowsForRow = getQueryReferencedRowsForRowImpl
	})

	t.Run("VerifyActualRowsWithDefect7", func(t *testing.T) {
		getGridForGridsApiImpl := getGridForGridsApi
		getGridForGridsApi = func(r ApiRequest, gridUuid string) (*model.Grid, error) { return nil, errors.New("xxx") } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"xxx"`)
		getGridForGridsApi = getGridForGridsApiImpl
	})

	t.Run("VerifyActualRowsWithDefect8", func(t *testing.T) {
		getGridQueryOutputForGridsApiImpl := getGridQueryOutputForGridsApi
		getGridQueryOutputForGridsApi = func(grid *model.Grid) []any { return nil } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when scanning grid definition: sql`)
		getGridQueryOutputForGridsApi = getGridQueryOutputForGridsApiImpl
	})

	t.Run("VerifyActualGridsOwnedByUser01Defect", func(t *testing.T) {
		getGridColumsNotOwnedQueryForGridsApiImpl := getGridColumsNotOwnedQueryForGridsApi
		getGridColumsNotOwnedQueryForGridsApi = func(bool) string { return "xxx" } // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when querying columns: pq: syntax error`)
		getGridColumsNotOwnedQueryForGridsApi = getGridColumsNotOwnedQueryForGridsApiImpl
	})

	t.Run("VerifyActualRowSingleDefect", func(t *testing.T) {
		getRowsQueryParametersForGridsApiImpl := getRowsQueryParametersForGridsApi
		getRowsQueryParametersForGridsApi = func(gridUuid, uuid string) []any { return nil }
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when querying rows: pq`)
		getRowsQueryParametersForGridsApi = getRowsQueryParametersForGridsApiImpl
	})

	t.Run("VerifyActualRowSingleDefect", func(t *testing.T) {
		getRowsQueryOutputForGridsApiImpl := getRowsQueryOutputForGridsApi
		getRowsQueryOutputForGridsApi = func(grid *model.Grid, row *model.Row) []any { return nil }
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when scanning rows: sql
		`)
		getRowsQueryOutputForGridsApi = getRowsQueryOutputForGridsApiImpl
	})
}
