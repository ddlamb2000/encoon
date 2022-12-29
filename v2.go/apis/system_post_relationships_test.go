// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"errors"
	"net/http"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestPostRelationships(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	var user01Uuid string
	db, _ := database.GetDbByName("test")
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test01").Scan(&user01Uuid)

	t.Run("CreateNewColumnsFor3rdGrid", func(t *testing.T) {
		var gridUuid1, gridUuid2 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid1)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid2)

		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 09","text2":"text1"},` +
			`{"uuid":"b","text1":"Test Column 10","text2":"text2"},` +
			`{"uuid":"c","text1":"Test Column 11","text2":"text3"},` +
			`{"uuid":"d","text1":"Test Column 12","text2":"text4"},` +
			`{"uuid":"e","text1":"Test Column 13","text2":"int1"},` +
			`{"uuid":"f","text1":"Test Column 14","text2":"int2"},` +
			`{"uuid":"g","text1":"Test Column 15","text2":"int3"},` +
			`{"uuid":"h","text1":"Test Column 16","text2":"int4"},` +
			`{"uuid":"i","text1":"Test Column 17","text2":"relationship1"},` +
			`{"uuid":"j","text1":"Test Column 18","text2":"relationship2"},` +
			`{"uuid":"k","text1":"Test Column 19","text2":"relationship3"},` +
			`{"uuid":"l","text1":"Test Column 20","text2":"relationship4"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"d","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"e","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"e","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"f","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"g","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"h","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"i","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"i","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid1 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"j","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"j","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid1 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"k","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"k","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid2 + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 09","text2":"text1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 10","text2":"text2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 11","text2":"text3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 12","text2":"text4"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 13","text2":"int1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 14","text2":"int2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 15","text2":"int3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 16","text2":"int4"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 17","text2":"relationship1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 18","text2":"relationship2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 19","text2":"relationship3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 20","text2":"relationship4"`)
		jsonStringContains(t, responseData, `"gridUuid":"`+model.UuidGrids+`","uuid":"`+gridUuid1+`"`)
		jsonStringContains(t, responseData, `"gridUuid":"`+model.UuidGrids+`","uuid":"`+gridUuid2+`"`)
	})

	t.Run("UpdateAndDeleteColumnRelationshipsFor3rdGrid", func(t *testing.T) {
		var gridUuid1, gridUuid2, column13Uuid, column19Uuid, column20Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid1)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid2)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 13").Scan(&column13Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 19").Scan(&column19Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 20").Scan(&column20Uuid)

		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"` + column19Uuid + `","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid1 + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + column20Uuid + `","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"` + column20Uuid + `","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid2 + `"}` +
			`],` +
			`"referencedValuesRemoved":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + column13Uuid + `","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"` + column19Uuid + `","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid2 + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 09","text2":"text1"`)
	})

	t.Run("Create3rdSingleGrid", func(t *testing.T) {
		var column09Uuid, column10Uuid, column11Uuid, column12Uuid, column13Uuid, column14Uuid, column15Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 09").Scan(&column09Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 10").Scan(&column10Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 11").Scan(&column11Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 12").Scan(&column12Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 13").Scan(&column13Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 14").Scan(&column14Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 15").Scan(&column15Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Grid03","text2":"Test grid 03","text3":"journal"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column09Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column10Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column11Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column12Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column13Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column14Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column15Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid03","text2":"Test grid 03","text3":"journal"`)
	})

	t.Run("AddColumnsTo3rdSingleGridDefect", func(t *testing.T) {
		getRowsQueryForGridUuidReferencedByColumnImpl := getRowsQueryForGridUuidReferencedByColumn
		getRowsQueryForGridUuidReferencedByColumn = func() string { return "xxx" }
		var gridUuid, column16Uuid, column17Uuid, column18Uuid, column19Uuid, column20Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 16").Scan(&column16Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 17").Scan(&column17Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 18").Scan(&column18Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 19").Scan(&column19Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 20").Scan(&column20Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column16Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column18Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column19Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column20Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		getRowsQueryForGridUuidReferencedByColumn = getRowsQueryForGridUuidReferencedByColumnImpl
	})

	t.Run("AddColumnsTo3rdSingleGridDefect2", func(t *testing.T) {
		removeAssociatedGridNotOwnedColumnFromCacheImpl := removeAssociatedGridNotOwnedColumnFromCache
		removeAssociatedGridNotOwnedColumnFromCache = func(apiRequestParameters, *model.Grid, string) error { return errors.New("xxx") }
		var gridUuid, column16Uuid, column17Uuid, column18Uuid, column19Uuid, column20Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 16").Scan(&column16Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 17").Scan(&column17Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 18").Scan(&column18Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 19").Scan(&column19Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 20").Scan(&column20Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column16Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column18Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column19Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column20Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		removeAssociatedGridNotOwnedColumnFromCache = removeAssociatedGridNotOwnedColumnFromCacheImpl
	})

	t.Run("AddColumnsTo3rdSingleGrid", func(t *testing.T) {
		var gridUuid, column16Uuid, column17Uuid, column18Uuid, column19Uuid, column20Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 16").Scan(&column16Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 17").Scan(&column17Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 18").Scan(&column18Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 19").Scan(&column19Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 20").Scan(&column20Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column16Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column18Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column19Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column20Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid03","text2":"Test grid 03","text3":"journal"`)
	})

	t.Run("VerifyNotOwnedColumnIn2ndSingleGrid", func(t *testing.T) {
		var grid2Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&grid2Uuid)
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid2Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":12`)
		jsonStringContains(t, responseData, `"owned":false,"label":"Test Column 20","name":"relationship4"`)
	})

	t.Run("CreateNewRowsIn3rdSingleGrid", func(t *testing.T) {
		var grid1Uuid, grid3Uuid, row01Uuid, row05Uuid, row17Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid1Uuid)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid1Uuid, "test-01").Scan(&row01Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid1Uuid, "test-05").Scan(&row05Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid1Uuid, "test-17").Scan(&row17Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a", "text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12","int1":13,"int2":14,"int3":15,"int4":15},` +
			`{"uuid":"b", "text1":"test-16","text2":"test-17","text3":"test-18","text4":"test-19","int1":20,"int2":21,"int3":22,"int4":23},` +
			`{"uuid":"c", "text1":"test-23","text2":"test-24","text3":"test-25","text4":"test-26","int1":27,"int2":28,"int3":29,"int4":30}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"owned":true,"columnName":"relationship3","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship3","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"owned":true,"columnName":"relationship3","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid3Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":3`)
		jsonStringContains(t, responseData, `"text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12","int1":13,"int2":14,"int3":15,"int4":15`)
		jsonStringContains(t, responseData, `"owned":true,"label":"Test Column 09","name":"text1","type":"Text"`)
		jsonStringContains(t, responseData, `"references":[{"owned":true,"label":"Test Column 17","name":"relationship1"`)
	})

	t.Run("CreateNewRowsIn3rdSingleGridDefect", func(t *testing.T) {
		getInsertStatementForReferenceRowImpl := getInsertStatementForReferenceRow
		getInsertStatementForReferenceRow = func() string { return "xxx" } // mock function
		var grid1Uuid, grid3Uuid, row05Uuid, row09Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid1Uuid)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid1Uuid, "test-05").Scan(&row05Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid3Uuid, "test-09").Scan(&row09Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"` + row09Uuid + `","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid3Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Insert referenced row error: pq: syntax error`)
		getInsertStatementForReferenceRow = getInsertStatementForReferenceRowImpl
	})

	t.Run("UpdateAndDeleteColumnRelationshipsFor3rdGridDefect", func(t *testing.T) {
		getDeleteReferenceRowStatementImpl := getDeleteReferenceRowStatement
		getDeleteReferenceRowStatement = func() string { return "xxx" } // mock function
		var grid1Uuid, grid3Uuid, row05Uuid, row09Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid1Uuid)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid1Uuid, "test-05").Scan(&row05Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid3Uuid, "test-09").Scan(&row09Uuid)

		postStr := `{"referencedValuesRemoved":` +
			`[` +
			`{"owned":true,"columnName":"relationship2","fromUuid":"` + row09Uuid + `","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid3Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusInternalServerError)
		jsonStringContains(t, responseData, `Delete referenced row error: pq: syntax error`)
		getDeleteReferenceRowStatement = getDeleteReferenceRowStatementImpl
	})

	t.Run("CreateNewRowIn3rdSingleGrid", func(t *testing.T) {
		var grid1Uuid, grid3Uuid, row01Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid1Uuid)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid1Uuid, "test-01").Scan(&row01Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a", "text1":"test-xx","text2":"test-yy","text3":"test-zz"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":true,"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid3Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":4`)
		jsonStringContains(t, responseData, `"text1":"test-xx","text2":"test-yy","text3":"test-zz"`)
		var rowXXUuid, referenceRowXXUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid3Uuid, "test-xx").Scan(&rowXXUuid)
		if rowXXUuid == "" {
			t.Errorf(`Row "test-xx" doesn't exist.`)
		}
		db.QueryRow("SELECT uuid FROM relationships WHERE gridUuid = $1 and text2= $2 and text3 = $3", model.UuidRelationships, grid3Uuid, rowXXUuid).Scan(&referenceRowXXUuid)
		if referenceRowXXUuid == "" {
			t.Errorf(`Referenced row for "test-xx" doesn't exist.`)
		}
	})

	t.Run("DeleteNewRowIn3rdSingleGrid", func(t *testing.T) {
		var grid3Uuid, rowXXUuid, referenceRowXXUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid3Uuid, "test-xx").Scan(&rowXXUuid)
		db.QueryRow("SELECT uuid FROM relationships WHERE gridUuid = $1 and text2= $2 and text3 = $3", model.UuidRelationships, grid3Uuid, rowXXUuid).Scan(&referenceRowXXUuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + rowXXUuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid3Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":3`)
		var newRowXXUuid, newReferenceRowXXUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2 and enabled = true", grid3Uuid, "test-xx").Scan(&newRowXXUuid)
		if newRowXXUuid != "" {
			t.Errorf(`Row "test-xx" still exists: %v.`, newRowXXUuid)
		}
		db.QueryRow("SELECT uuid FROM relationships WHERE gridUuid = $1 and text2= $2 and text3 = $3  and enabled = true", model.UuidRelationships, grid3Uuid, rowXXUuid).Scan(&newReferenceRowXXUuid)
		if newReferenceRowXXUuid != "" {
			t.Errorf(`Referenced row for "test-xx" still exists: %v.`, newReferenceRowXXUuid)
		}
	})

	t.Run("CreateNewRowsIn2ndSingleGridWithNotOwnedReferences", func(t *testing.T) {
		var grid2Uuid, grid3Uuid, row09Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&grid2Uuid)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid3Uuid, "test-09").Scan(&row09Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","int1":2,"int2":3,"int3":4,"int4":5},` +
			`{"uuid":"b","int1":3,"int2":4,"int3":5,"int4":6}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"owned":false,"columnName":"relationship4","fromUuid":"a","toGridUuid":"` + grid3Uuid + `","uuid":"` + row09Uuid + `"},` +
			`{"owned":false,"columnName":"relationship4","fromUuid":"b","toGridUuid":"` + grid3Uuid + `","uuid":"` + row09Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid2Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":14`)
		jsonStringContains(t, responseData, `"int1":2,"int2":3,"int3":4,"int4":5`)
		jsonStringContains(t, responseData, `"int1":3,"int2":4,"int3":5,"int4":6`)
		jsonStringContains(t, responseData, `"owned":false,"label":"Test Column 20","name":"relationship4"`)
		jsonStringContains(t, responseData, `"text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12"`)
	})

	t.Run("RemoveNotOwnedReferencesIn2ndSingleGrid", func(t *testing.T) {
		var grid2Uuid, grid3Uuid, row09Uuid, rowInt3Uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&grid2Uuid)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid3Uuid, "test-09").Scan(&row09Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and int1= $2", grid2Uuid, 3).Scan(&rowInt3Uuid)
		postStr := `{"referencedValuesRemoved":` +
			`[` +
			`{"owned":false,"columnName":"relationship4","fromUuid":"` + rowInt3Uuid + `","toGridUuid":"` + grid3Uuid + `","uuid":"` + row09Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid2Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"int1":3,"int2":4,"int3":5,"int4":6`)
	})
}
