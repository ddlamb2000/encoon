// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
)

func RunSystemTestPost(t *testing.T) {
	t.Run("CreateNewUserNoAuth", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"aaaa","text02":"bbbb"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("te st", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringContains(t, responseData, `"error":"Invalid request or unauthorized database access: signature is invalid."`)
	})

	t.Run("Post404", func(t *testing.T) {
		postStr := `{}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		byteEqualString(t, responseData, `404 page not found`)
	})

	t.Run("CreateUserNoData", func(t *testing.T) {
		postStr := `{}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+utils.UuidUsers+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `{"uuid":"`+utils.UuidRootUser+`"`)
	})

	t.Run("CreateNewSingleUser", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test01","text02":"Zero-one","text03":"Test","text04":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringDoesntContain(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"countRows":2`)
		jsonStringContains(t, responseData, `"grid":{"uuid":"`+utils.UuidUsers+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `{"uuid":"`+utils.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"Zero-one","text03":"Test"`)
	})

	t.Run("Create3Users", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test02","text02":"Zero-two","text03":"Test","text04":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"},` +
			`{"text01":"test03","text02":"Zero-three","text03":"Test","text04":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"},` +
			`{"text01":"test04","text02":"Zero-four","text03":"Test","text04":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringDoesntContain(t, responseData, `"countRows":2`)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringContains(t, responseData, `"text01":"test02","text02":"Zero-two","text03":"Test"`)
		jsonStringContains(t, responseData, `"text01":"test03","text02":"Zero-three","text03":"Test"`)
		jsonStringContains(t, responseData, `"text01":"test04","text02":"Zero-four","text03":"Test"`)
	})

	t.Run("CreateNewSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"Grid01","text02":"Test grid 01","text03":"Grid generated by system test.","text04":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":7`)
		jsonStringContains(t, responseData, `"text01":"Grid01","text02":"Test grid 01","text03":"Grid generated by system test.","text04":"journal"`)
	})

	t.Run("VerifySingleGridIsCreated", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":7`)
		jsonStringContains(t, responseData, `"text01":"Grid01","text02":"Test grid 01","text03":"Grid generated by system test.","text04":"journal"`)
	})

	t.Run("VerifyNoRowInSingleGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":0`)
		jsonStringContains(t, responseData, `"rows":[]`)
	})

	t.Run("CreateNewRowInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test01","text02":"test02","text03":"test03","text04":"test04"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"test02","text03":"test03","text04":"test04"`)
	})

	t.Run("CreateRowIncorrectPayload", func(t *testing.T) {
		postStr := `{"xxxxx"}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"test02","text03":"test03","text04":"test04"`)
	})

	t.Run("VerifyNewRowInSingleGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"test02","text03":"test03","text04":"test04"`)
	})

	t.Run("CreateNewRowInSingleGridWithTimeOut", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test05","text02":"test06","text03":"test07","text04":"test08"}` +
			`]` +
			`}`
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusRequestTimeout)
		jsonStringContains(t, responseData, `{"error":"[test] [root] Post request has been cancelled: context deadline exceeded."}`)
	})

	t.Run("VerifyNoNewRowInSingleGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"test02","text03":"test03","text04":"test04"`)
		jsonStringDoesntContain(t, responseData, `"text01":"test05"`)
	})

	t.Run("UpdateNewRow", func(t *testing.T) {
		db := database.GetDbByName("test")
		var uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", utils.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuid + `","text01":"Grid01","text02":"Test grid 01","text03":"Grid 01 generated by system test.","text04":"journal"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"uuid":"`+uuid+`"`)
		jsonStringContains(t, responseData, `"text01":"Grid01","text02":"Test grid 01","text03":"Grid 01 generated by system test.","text04":"journal"`)
		var version int
		db.QueryRow("SELECT version FROM rows WHERE gridUuid = $1 and uuid = $2", utils.UuidGrids, uuid).Scan(&version)
		intEqual(t, version, 2)
	})

	t.Run("CreateNewandUpdateRowsInSingleGrid", func(t *testing.T) {
		db := database.GetDbByName("test")
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", utils.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", uuidGrid, "test01").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test-05","text02":"test-06","text03":"test-07","text04":"test-08"},` +
			`{"text01":"test-09","text02":"test-10","text03":"test-11","text04":"test-12"},` +
			`{"text01":"test-13","text02":"test-14","text03":"test-15","text04":"test-16"},` +
			`{"text01":"test-17","text02":"test-18","text03":"test-19","text04":"test-20"}` +
			`],` +
			`"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuidRow + `","text01":"test-01","text02":"test-02","text03":"test-03","text04":"test-04"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringContains(t, responseData, `"text01":"test-01","text02":"test-02","text03":"test-03","text04":"test-04"`)
		jsonStringContains(t, responseData, `"text01":"test-09","text02":"test-10","text03":"test-11","text04":"test-12"`)
	})

	t.Run("CreateDeleteRowsInSingleGrid", func(t *testing.T) {
		db := database.GetDbByName("test")
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", utils.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", uuidGrid, "test-13").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test-21","text02":"test-22","text03":"test-23","text04":"test-24"},` +
			`{"text01":"test-25","text02":"test-26","text03":"test-27","text04":"test-28"}` +
			`],` +
			`"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text01":"test-21","text02":"test-22","text03":"test-23","text04":"test-24"`)
		jsonStringContains(t, responseData, `"text01":"test-25","text02":"test-26","text03":"test-27","text04":"test-28"`)
		jsonStringDoesntContain(t, responseData, `"text01":"test-13"`)
	})

	t.Run("UpdateIncorrectRowInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"xxxx","text01":"test-01","text02":"test-02","text03":"test-03","text04":"test-04"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringContains(t, responseData, `"error":"[test] [root] Update row error`)
	})

	t.Run("DeleteIncorrectRowInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"xxxx"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		jsonStringContains(t, responseData, `"error":"[test] [root] Delete row error`)
	})

	t.Run("DeleteRowInSingleGrid", func(t *testing.T) {
		db := database.GetDbByName("test")
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", utils.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text01= $2", uuidGrid, "test-09").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringDoesntContain(t, responseData, `"text01":"test-09"`)
	})

	t.Run("CreateNewRowInSingleGridWithTimeOut2", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"test-29","text02":"test-30","text03":"test-31","text04":"test-32"}` +
			`]` +
			`}`
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 10, 500)
		responseData, code, err := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01?trace=true", postStr)
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text01":"test-29","text02":"test-30","text03":"test-31","text04":"test-32"`)
	})
}
