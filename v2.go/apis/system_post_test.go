// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestPost(t *testing.T) {
	t.Run("CreateNewSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":8`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"journal"`)
	})

	t.Run("VerifySingleGridIsCreated", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":8`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"journal"`)
	})

	t.Run("VerifyNoRowInSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":0`)
		jsonStringContains(t, responseData, `"rows":[]`)
	})

	t.Run("CreateNewColumnsInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Test Column 01","text2":"text1"},` +
			`{"text1":"Test Column 02","text2":"text2"},` +
			`{"text1":"Test Column 03","text2":"text3"},` +
			`{"text1":"Test Column 04","text2":"text4"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 01","text2":"text1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 02","text2":"text2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 03","text2":"text3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 04","text2":"text4"`)

		db, _ := database.GetDbByName("test")
		var gridUuid, uuidCol1, uuidCol2, uuidCol3, uuidCol4 string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 01").Scan(&uuidCol1)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 02").Scan(&uuidCol2)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 03").Scan(&uuidCol3)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 04").Scan(&uuidCol4)

		postStr = `{"rowsAdded":` +
			`[` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol1 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol2 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol3 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol4 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol1 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidTextColumnType + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol2 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidTextColumnType + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol3 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidTextColumnType + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol4 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidTextColumnType + `"}` +
			`]` +
			`}`
		responseData, code, err = runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidRelationships, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
	})

	t.Run("CreateNewRowInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test01","text2":"test02","text3":"test03","text4":"test04"}` +
			`]` +
			`}`
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Test Column 01","name":"text1","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 02","name":"text2","type":"Text"`)
	})

	t.Run("CreateRowIncorrectPayload", func(t *testing.T) {
		postStr := `{"xxxxx"}`
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
	})

	t.Run("VerifyNewRowInSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
	})

	t.Run("CreateNewRowInSingleGridWithTimeOut", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test05","text2":"test06","text3":"test07","text4":"test08"}` +
			`]` +
			`}`
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		defer database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusRequestTimeout)
		jsonStringContains(t, responseData, `{"error":"Post request has been cancelled: context deadline exceeded."}`)
	})

	t.Run("VerifyNoNewRowInSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
		jsonStringDoesntContain(t, responseData, `"text1":"test05"`)
	})

	t.Run("UpdateNewRow", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuid + `","text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"uuid":"`+uuid+`"`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"journal"`)
		var version int
		db.QueryRow("SELECT version FROM rows WHERE gridUuid = $1 and uuid = $2", model.UuidGrids, uuid).Scan(&version)
		intEqual(t, version, 2)
	})

	t.Run("CreateNewandUpdateRowsInSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test01").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-05","text2":"test-06","text3":"test-07","text4":"test-08"},` +
			`{"text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12"},` +
			`{"text1":"test-13","text2":"test-14","text3":"test-15","text4":"test-16"},` +
			`{"text1":"test-17","text2":"test-18","text3":"test-19","text4":"test-20"}` +
			`],` +
			`"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuidRow + `","text1":"test-01","text2":"test-02","text3":"test-03","text4":"test-04"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+uuidGrid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringContains(t, responseData, `"text1":"test-01","text2":"test-02","text3":"test-03","text4":"test-04"`)
		jsonStringContains(t, responseData, `"text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12"`)
	})

	t.Run("CreateDeleteRowsInSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-13").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-21","text2":"test-22","text3":"test-23","text4":"test-24"},` +
			`{"text1":"test-25","text2":"test-26","text3":"test-27","text4":"test-28"}` +
			`],` +
			`"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+uuidGrid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text1":"test-21","text2":"test-22","text3":"test-23","text4":"test-24"`)
		jsonStringContains(t, responseData, `"text1":"test-25","text2":"test-26","text3":"test-27","text4":"test-28"`)
		jsonStringDoesntContain(t, responseData, `"text1":"test-13"`)
	})

	t.Run("DeleteRowInSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-09").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+uuidGrid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringDoesntContain(t, responseData, `"text1":"test-09"`)
	})

	t.Run("CreateNewRowInSingleGridWithTimeOut2", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-29","text2":"test-30","text3":"test-31","text4":"test-32"}` +
			`]` +
			`}`
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 10, 500)
		defer database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text1":"test-29","text2":"test-30","text3":"test-31","text4":"test-32"`)
	})

	t.Run("Create2ndSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid02","text2":"Test grid 02","text3":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":9`)
		jsonStringContains(t, responseData, `"text1":"Grid02","text2":"Test grid 02","text3":"journal"`)
	})

	t.Run("CreateNewColumnsIn2ndGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Test Column 05","text2":"int1"},` +
			`{"text1":"Test Column 06","text2":"int2"},` +
			`{"text1":"Test Column 07","text2":"int3"},` +
			`{"text1":"Test Column 08","text2":"int4"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 05","text2":"int1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 06","text2":"int2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 07","text2":"int3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 08","text2":"int4"`)

		db, _ := database.GetDbByName("test")
		var gridUuid, uuidCol1, uuidCol2, uuidCol3, uuidCol4 string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 05").Scan(&uuidCol1)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 06").Scan(&uuidCol2)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 07").Scan(&uuidCol3)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 08").Scan(&uuidCol4)

		postStr = `{"rowsAdded":` +
			`[` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol1 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol2 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol3 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidGrids + `", "text3":"` + gridUuid + `", "text4":"` + model.UuidColumns + `", "text5":"` + uuidCol4 + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol1 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidIntColumnType + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol2 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidIntColumnType + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol3 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidIntColumnType + `"},` +
			`{"text1":"relationship1","text2":"` + model.UuidColumns + `", "text3":"` + uuidCol4 + `", "text4":"` + model.UuidColumnTypes + `", "text5":"` + model.UuidIntColumnType + `"}` +
			`]` +
			`}`
		responseData, code, err = runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidRelationships, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
	})

	t.Run("CreateNewRowsIn2ndSingleGrid", func(t *testing.T) {
		db, _ := database.GetDbByName("test")
		var gridUuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+gridUuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"countRows":12`)
		jsonStringContains(t, responseData, `"int1":1,"int2":2,"int3":3,"int4":4`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Test Column 05","name":"int1","type":"Integer"`)
	})
}
