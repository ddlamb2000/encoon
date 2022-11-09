// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/utils"
)

func RunSystemTestPost(t *testing.T) {
	t.Run("CreateNewUserNoAuth", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text01":"aaaa","text02":"bbbb"}` +
			`]` +
			`}`
		responseData, err, code := runPOSTRequestForUser("te st", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusUnauthorized)
		jsonStringContains(t, responseData, `"error":"Invalid request or unauthorized database access: signature is invalid."`)
	})

	t.Run("Post404", func(t *testing.T) {
		postStr := `{}`
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusNotFound)
		byteEqualString(t, responseData, `404 page not found`)
	})

	t.Run("CreateUserNoData", func(t *testing.T) {
		postStr := `{}`
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
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
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
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
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_users", postStr)
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
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":7`)
		jsonStringContains(t, responseData, `"text01":"Grid01","text02":"Test grid 01","text03":"Grid generated by system test.","text04":"journal"`)
	})

	t.Run("VerifySingleGridIsCreated", func(t *testing.T) {
		responseData, err, code := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/_grids")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":7`)
		jsonStringContains(t, responseData, `"text01":"Grid01","text02":"Test grid 01","text03":"Grid generated by system test.","text04":"journal"`)
	})

	t.Run("VerifyNoRowInSingleGrid", func(t *testing.T) {
		responseData, err, code := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01")
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
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"test02","text03":"test03","text04":"test04"`)
	})

	t.Run("VerifyNewRowInSingleGrid", func(t *testing.T) {
		responseData, err, code := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01")
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
		forceTestSleepTime("test", 500)
		forceTimeOutThreshold(200)
		responseData, err, code := runPOSTRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01", postStr)
		forceTestSleepTime("test", 0)
		forceTimeOutThreshold(200)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusRequestTimeout)
		jsonStringContains(t, responseData, `{"error":"[test] [root] Post request has been cancelled: context deadline exceeded."}`)
	})

	t.Run("VerifyNoNewRowInSingleGrid", func(t *testing.T) {
		responseData, err, code := runGETRequestForUser("test", "root", utils.UuidRootUser, "/test/api/v1/Grid01")
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text01":"test01","text02":"test02","text03":"test03","text04":"test04"`)
		jsonStringDoesntContain(t, responseData, `"text01":"test05"`)
	})

}
