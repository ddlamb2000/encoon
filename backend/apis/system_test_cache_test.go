// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"errors"
	"net/http"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestCache(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	var user01Uuid, user02Uuid, user03Uuid, grid01Uuid, grid02Uuid, grid03Uuid string
	db, _ := database.GetDbByName("test")
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test01").Scan(&user01Uuid)
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test02").Scan(&user02Uuid)
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test03").Scan(&user03Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid01Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&grid02Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid03Uuid)
	var row17Uuid, rowInt100Uuid, row23Uuid string
	db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid01Uuid, "test-17").Scan(&row17Uuid)
	db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and int1= $2", grid02Uuid, 100).Scan(&rowInt100Uuid)
	db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid03Uuid, "test-23").Scan(&row23Uuid)

	t.Run("User01VerifyActualGridsCount", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, requestContent{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":10`)
	})

	t.Run("User01CreateNewColumnsFor5thGrid", func(t *testing.T) {
		var gridUuid3 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)

		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 24","text2":"text1"},` +
			`{"uuid":"b","text1":"Test Column 25","text2":"text2"},` +
			`{"uuid":"c","text1":"Test Column 26","text2":"relationship1","text3":"true"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"c","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid3 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 24","text2":"text1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 25","text2":"text2"`)
	})

	t.Run("User01Create5thSingleGrid", func(t *testing.T) {
		var column24Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 24").Scan(&column24Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Grid05","text2":"Test grid 05","text3":"journal"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column24Uuid + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"a","toGridUuid":"` + model.UuidAccessLevels + `","uuid":"` + model.UuidAccessLevelWriteAccess + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid05","text2":"Test grid 05","text3":"journal"`)
	})

	t.Run("User01VerifyActualGridsCount2", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, requestContent{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":11`)
	})

	t.Run("User01VerifyActualGridsColumns", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, grid05Uuid, requestContent{
			Action:   ActionLoad,
			GridUuid: grid05Uuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"label":"Test Column 24","name":"text1","type":"Text"`)
	})

	t.Run("User01AddColumnsTo5thSingleGridDefect", func(t *testing.T) {
		removeAssociatedGridFromCacheImpl := removeAssociatedGridFromCache
		removeAssociatedGridFromCache = func(ApiRequest, *model.Grid, string) error { return errors.New("xxx") }
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		var column25Uuid, column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 25").Scan(&column25Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26").Scan(&column26Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + grid05Uuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column25Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + grid05Uuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column26Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		removeAssociatedGridFromCache = removeAssociatedGridFromCacheImpl
	})

	t.Run("User01AddColumnsTo5thSingleGrid", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		var column25Uuid, column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 25").Scan(&column25Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26").Scan(&column26Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + grid05Uuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column25Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + grid05Uuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column26Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid05","text2":"Test grid 05","text3":"journal"`)
	})

	t.Run("User01VerifyActualGridsColumns2", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, grid05Uuid, requestContent{
			Action:   ActionLoad,
			GridUuid: grid05Uuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"label":"Test Column 24","name":"text1","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 25","name":"text2","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 26","name":"relationship1","type":"Reference"`)
	})

	t.Run("User01RemoveColumnFrom5thSingleGridDefect", func(t *testing.T) {
		removeAssociatedGridFromCacheImpl := removeAssociatedGridFromCache
		removeAssociatedGridFromCache = func(ApiRequest, *model.Grid, string) error { return errors.New("xxx") }
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		var column25Uuid, column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 25").Scan(&column25Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26").Scan(&column26Uuid)
		postStr := `{"referencedValuesRemoved":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + grid05Uuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column25Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		removeAssociatedGridFromCache = removeAssociatedGridFromCacheImpl
	})

	t.Run("User01RemoveColumnFrom5thSingleGrid", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		var column25Uuid, column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 25").Scan(&column25Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26").Scan(&column26Uuid)
		postStr := `{"referencedValuesRemoved":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + grid05Uuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column25Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid05","text2":"Test grid 05","text3":"journal"`)
	})

	t.Run("User01VerifyActualGridsColumns3", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, grid05Uuid, requestContent{
			Action:   ActionLoad,
			GridUuid: grid05Uuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"label":"Test Column 24","name":"text1","type":"Text"`)
		jsonStringDoesntContain(t, responseData, `"label":"Test Column 25","name":"text2","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 26","name":"relationship1","type":"Reference"`)
	})

	t.Run("User01Rename5thSingleGridDefect", func(t *testing.T) {
		removeAssociatedGridFromCacheImpl := removeAssociatedGridFromCache
		removeAssociatedGridFromCache = func(ApiRequest, *model.Grid, string) error { return errors.New("xxx") }
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + grid05Uuid + `","text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		removeAssociatedGridFromCache = removeAssociatedGridFromCacheImpl
	})

	t.Run("User01Rename5thSingleGrid", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05").Scan(&grid05Uuid)
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + grid05Uuid + `","text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"`)
	})

	t.Run("User01VerifyActualGridName", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05 {2}").Scan(&grid05Uuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, grid05Uuid, requestContent{
			Action:   ActionLoad,
			GridUuid: grid05Uuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"`)
		jsonStringContains(t, responseData, `"label":"Test Column 24","name":"text1","type":"Text"`)
		jsonStringDoesntContain(t, responseData, `"label":"Test Column 25","name":"text2","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 26","name":"relationship1","type":"Reference"`)
	})

	t.Run("User01RenameColumnFor5thSingleGridDefect", func(t *testing.T) {
		getGridUuidAttachedToColumnImpl := getGridUuidAttachedToColumn
		getGridUuidAttachedToColumn = func(ApiRequest, string) (string, error) { return "", errors.New("xxx") }
		var column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26").Scan(&column26Uuid)
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + column26Uuid + `","text1":"Test Column 26 {2}","text2":"relationship1","text3":"true"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		getGridUuidAttachedToColumn = getGridUuidAttachedToColumnImpl
	})

	t.Run("User01RenameColumnFor5thSingleGrid", func(t *testing.T) {
		var column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26").Scan(&column26Uuid)
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + column26Uuid + `","text1":"Test Column 26 {2}","text2":"relationship1","text3":"true"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 26 {2}"`)
	})

	t.Run("User01VerifyActualColumnName", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05 {2}").Scan(&grid05Uuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, grid05Uuid, requestContent{
			Action:   ActionLoad,
			GridUuid: grid05Uuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"`)
		jsonStringContains(t, responseData, `"label":"Test Column 24","name":"text1","type":"Text"`)
		jsonStringDoesntContain(t, responseData, `"label":"Test Column 25","name":"text2","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 26 {2}","name":"relationship1","type":"Reference"`)
	})

	t.Run("User01Remove5thSingleGrid", func(t *testing.T) {
		var grid05Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid05 {2}").Scan(&grid05Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + grid05Uuid + `","text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)

		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, requestContent{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringDoesntContain(t, responseData, `"text1":"Grid05 {2}","text2":"Test grid 05 {2}","text3":"journal"`)
	})

	t.Run("User01VerifyActualGridsCount2", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, requestContent{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":10`)
	})

	t.Run("User01DeleteColumnFor5thSingleGridDefect", func(t *testing.T) {
		getGridUuidAttachedToColumnForCacheImpl := getGridUuidAttachedToColumnForCache
		getGridUuidAttachedToColumnForCache = func(ApiRequest, string) (string, error) { return "", errors.New("xxx") }
		var column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26 {2}").Scan(&column26Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + column26Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `"error":"Error when getting data for cache deletion: xxx."`)
		getGridUuidAttachedToColumnForCache = getGridUuidAttachedToColumnForCacheImpl
	})

	t.Run("User01DeleteColumnFor5thSingleGrid", func(t *testing.T) {
		var column26Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 26 {2}").Scan(&column26Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + column26Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
	})
}
