// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"errors"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestNotOwnedColumn(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	var user01Uuid string
	db, _ := database.GetDbByName("test")
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test01").Scan(&user01Uuid)

	t.Run("User01Create6thSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Grid06","text2":"Test grid 06","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, requestContent{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Grid06","text2":"Test grid 06","text3":"journal"`)
	})

	t.Run("User01CreateNewColumnsFor6thGrid", func(t *testing.T) {
		var gridUuid3, gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)

		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 27","text2":"text1"},` +
			`{"uuid":"b","text1":"Test Column 28","text2":"int1"},` +
			`{"uuid":"c","text1":"Test Column 29","text2":"text2"},` +
			`{"uuid":"d","text1":"Test Column 30","text2":"relationship1","text3":"true"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidBooleanColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"d","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid3 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"d","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"d","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidColumns,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Test Column 27","text2":"text1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 28","text2":"int1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 29","text2":"text2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 30","text2":"relationship1"`)
	})

	t.Run("CreateNewRowIn6thGrid", func(t *testing.T) {
		var gridUuid3, gridUuid6, row09Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", gridUuid3, "test-09").Scan(&row09Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"test01","int1":1,"text2":""},` +
			`{"uuid":"b","text1":"test02","int1":2,"text2":"true"},` +
			`{"uuid":"c","text1":"test03","int1":3,"text2":""},` +
			`{"uuid":"d","text1":"test04","int1":4,"text2":"true"},` +
			`{"uuid":"e","text1":"test05","int1":5,"text2":""},` +
			`{"uuid":"f","text1":"test06","int1":6,"text2":"true"},` +
			`{"uuid":"g","text1":"test07","int1":7,"text2":""},` +
			`{"uuid":"h","text1":"test08","int1":8,"text2":"true"},` +
			`{"uuid":"i","text1":"test09","int1":9,"text2":""},` +
			`{"uuid":"j","text1":"test10","int1":10,"text2":"true"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"d","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"e","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"f","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"g","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"h","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"i","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"j","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid6, requestContent{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid6,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":10`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"","int1":1`)
		jsonStringContains(t, responseData, `"text1":"test02","text2":"true","int1":2`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Test Column 27","name":"text1","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 28","name":"int1","type":"Integer"`)
		jsonStringContains(t, responseData, `"label":"Test Column 29","name":"text2","type":"Boolean"`)
		jsonStringContains(t, responseData, `"label":"Test Column 30","name":"relationship1","type":"Reference"`)
		jsonStringContains(t, responseData, `"displayString":"test-09"`)
	})

	t.Run("VerifyActualGridsOwnedByUser01Count", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, requestContent{
			Action:               ActionLoad,
			GridUuid:             model.UuidGrids,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship3",
			FilterColumnGridUuid: model.UuidUsers,
			FilterColumnValue:    user01Uuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":5`)
	})

	t.Run("User01CreateAdditionalColumnsWithDefaultFor6thGridDefect", func(t *testing.T) {
		getInsertStatementForReferenceRowImpl := getInsertStatementForReferenceRow
		getInsertStatementForReferenceRow = func() string { return "xxx" } // mock function
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 31","text2":"text3"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Insert referenced row error: pq: syntax error`)
		getInsertStatementForReferenceRow = getInsertStatementForReferenceRowImpl
	})

	t.Run("User01CreateAdditionalColumnsWithDefaultFor6thGrid", func(t *testing.T) {
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 31"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Test Column 31","text2":"text3"`)
	})

	t.Run("CreateAdditionalRowsIn6thGrid", func(t *testing.T) {
		var gridUuid3, gridUuid6, row09Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", gridUuid3, "test-09").Scan(&row09Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"test11","int1":11,"text2":"","text3":"12"},` +
			`{"uuid":"b","text1":"test13","int1":13,"text2":"","text3":"14"},` +
			`{"uuid":"c","text1":"test15","int1":15,"text2":"","text3":"16"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid6, requestContent{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid6,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":13`)
		jsonStringContains(t, responseData, `"text1":"test11","text2":"","text3":"12","int1":11`)
		jsonStringContains(t, responseData, `"text1":"test13","text2":"","text3":"14","int1":13`)
		jsonStringContains(t, responseData, `"text1":"test15","text2":"","text3":"16","int1":15`)
		jsonStringContains(t, responseData, `"owned":true,"label":"Test Column 31","name":"text3","type":"Text"`)
	})

	t.Run("User01CreateAnotherColumnsWithDefaultFor6thGridDefect1", func(t *testing.T) {
		getUpdateColumnOrderNumberQueryImpl := getUpdateColumnOrderNumberQuery
		getUpdateColumnOrderNumberQuery = func() string { return "xxx" }
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 32"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, _ := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsFailure(t, response)
		getUpdateColumnOrderNumberQuery = getUpdateColumnOrderNumberQueryImpl
	})

	t.Run("User01CreateAnotherColumnsWithDefaultFor6thGridDefect2", func(t *testing.T) {
		getUpdateColumnNameQueryImpl := getUpdateColumnNameQuery
		getUpdateColumnNameQuery = func() string { return "xxx" }
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 32"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, _ := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsFailure(t, response)
		getUpdateColumnNameQuery = getUpdateColumnNameQueryImpl
	})

	t.Run("User01CreateAnotherColumnsWithDefaultFor6thGridDefect3", func(t *testing.T) {
		getGridInstanceWithColumnsForGridsApiImpl := getGridInstanceWithColumnsForGridsApi
		getGridInstanceWithColumnsForGridsApi = func(ApiRequest, string) (*model.Grid, error) { return nil, errors.New("xxx") }
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 32"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, _ := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsFailure(t, response)
		getGridInstanceWithColumnsForGridsApi = getGridInstanceWithColumnsForGridsApiImpl
	})

	t.Run("User01CreateAnotherColumnsWithDefaultFor6thGridDefect4", func(t *testing.T) {
		getGridUuidAttachedToColumnToUpdateWithColumnDefaultsImpl := getGridUuidAttachedToColumnToUpdateWithColumnDefaults
		getGridUuidAttachedToColumnToUpdateWithColumnDefaults = func(ApiRequest, string) (string, error) { return "", errors.New("xxx") }
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 32"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, _ := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		getGridUuidAttachedToColumnToUpdateWithColumnDefaults = getGridUuidAttachedToColumnToUpdateWithColumnDefaultsImpl
	})

	t.Run("User01CreateAnotherColumnsWithDefaultFor6thGrid", func(t *testing.T) {
		var gridUuid6 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid06").Scan(&gridUuid6)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 32"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":false,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid6 + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, requestContent{
			Action:               ActionChangeGrid,
			GridUuid:             model.UuidColumns,
			FilterColumnOwned:    true,
			FilterColumnName:     "relationship1",
			FilterColumnGridUuid: model.UuidColumnTypes,
			FilterColumnValue:    model.UuidTextColumnType,
			DataSet:              stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Test Column 32","text2":"text4"`)
	})
}
