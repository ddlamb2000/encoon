package apis

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestPostRelationships(t *testing.T) {
	t.Run("CreateNewColumnsFor3rdGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid1, gridUuid2 string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid1)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid2)

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
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"d","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"e","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"e","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"f","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"g","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"h","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidIntColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"i","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"columnName":"relationship2","fromUuid":"i","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid1 + `"},` +
			`{"columnName":"relationship1","fromUuid":"j","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"columnName":"relationship2","fromUuid":"j","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid1 + `"},` +
			`{"columnName":"relationship1","fromUuid":"k","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"columnName":"relationship2","fromUuid":"k","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid2 + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidColumns, postStr)
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
		db, _ := database.GetDbByName("test")
		var gridUuid1, gridUuid2, column13Uuid, column19Uuid, column20Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid1)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid2)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 13").Scan(&column13Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 19").Scan(&column19Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 20").Scan(&column20Uuid)

		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship2","fromUuid":"` + column19Uuid + `","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid1 + `"},` +
			`{"columnName":"relationship1","fromUuid":"` + column20Uuid + `","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"},` +
			`{"columnName":"relationship2","fromUuid":"` + column20Uuid + `","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid2 + `"}` +
			`],` +
			`"referencedValuesRemoved":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"` + column13Uuid + `","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship2","fromUuid":"` + column19Uuid + `","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid2 + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 09","text2":"text1"`)
	})

	t.Run("Create3rdSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var column09Uuid, column10Uuid, column11Uuid, column12Uuid, column13Uuid, column14Uuid, column15Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 09").Scan(&column09Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 10").Scan(&column10Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 11").Scan(&column11Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 12").Scan(&column12Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 13").Scan(&column13Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 14").Scan(&column14Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 15").Scan(&column15Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Grid03","text2":"Test grid 03","text3":"journal"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column09Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column10Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column11Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column12Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column13Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column14Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column15Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":10`)
		jsonStringContains(t, responseData, `"text1":"Grid03","text2":"Test grid 03","text3":"journal"`)
	})

	t.Run("AddColumnsTo3rdSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid, column16Uuid, column17Uuid, column18Uuid, column19Uuid, column20Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 16").Scan(&column16Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 17").Scan(&column17Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 18").Scan(&column18Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 19").Scan(&column19Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 20").Scan(&column20Uuid)
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column16Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column17Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column18Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column19Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"` + gridUuid + `","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column20Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":10`)
		jsonStringContains(t, responseData, `"text1":"Grid03","text2":"Test grid 03","text3":"journal"`)
	})

	t.Run("CreateNewRowsIn3rdSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var grid1Uuid, grid3Uuid, row01Uuid, row05Uuid, row17Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid1Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid3Uuid)
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
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"},` +
			`{"columnName":"relationship2","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"columnName":"relationship3","fromUuid":"a","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"},` +
			`{"columnName":"relationship2","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"columnName":"relationship2","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"columnName":"relationship3","fromUuid":"b","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row01Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"},` +
			`{"columnName":"relationship2","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row05Uuid + `"},` +
			`{"columnName":"relationship3","fromUuid":"c","toGridUuid":"` + grid1Uuid + `","uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+grid3Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":3`)
		jsonStringContains(t, responseData, `"text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12","int1":13,"int2":14,"int3":15,"int4":15`)
		jsonStringContains(t, responseData, `"columns":[{"label":"Test Column 09","name":"text1","type":"Text"`)
		jsonStringContains(t, responseData, `"references":[{"label":"Test Column 17","name":"relationship1","rows":[{`)
	})
}
