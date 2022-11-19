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

	t.Run("CreateNewSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid01","text2":"Test grid 01","text3":"Grid generated by system test.","text4":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_grids", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":8`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"Grid generated by system test.","text4":"journal"`)
	})

	t.Run("VerifySingleGridIsCreated", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_grids")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":8`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"Grid generated by system test.","text4":"journal"`)
	})

	t.Run("VerifyNoRowInSingleGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01")
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
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_columns", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
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
		responseData, code, err = runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_relationships", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
	})

	t.Run("CreateNewRowInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test01","text2":"test02","text3":"test03","text4":"test04"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
	})

	t.Run("CreateRowIncorrectPayload", func(t *testing.T) {
		postStr := `{"xxxxx"}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
	})

	t.Run("VerifyNewRowInSingleGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01")
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
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusRequestTimeout)
		jsonStringContains(t, responseData, `{"error":"[test] [root] Post request has been cancelled: context deadline exceeded."}`)
	})

	t.Run("VerifyNoNewRowInSingleGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01")
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
			`{"uuid":"` + uuid + `","text1":"Grid01","text2":"Test grid 01","text3":"Grid 01 generated by system test.","text4":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/_grids", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"uuid":"`+uuid+`"`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"Grid 01 generated by system test.","text4":"journal"`)
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
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
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
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
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
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringDoesntContain(t, responseData, `"text1":"test-09"`)
	})

	t.Run("CreateNewRowInSingleGridWithTimeOut2", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-29","text2":"test-30","text3":"test-31","text4":"test-32"}` +
			`]` +
			`}`
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 10, 500)
		defer database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		responseData, code, err := runPOSTRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/Grid01?trace=true", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text1":"test-29","text2":"test-30","text3":"test-31","text4":"test-32"`)
	})
}
