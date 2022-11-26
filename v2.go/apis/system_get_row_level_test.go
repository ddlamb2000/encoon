// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"net/http"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestGetRowLevel(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	var user01Uuid, user02Uuid, user03Uuid, grid01Uuid, grid02Uuid, grid03Uuid string
	db, _ := database.GetDbByName("test")
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test01").Scan(&user01Uuid)
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test02").Scan(&user02Uuid)
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test03").Scan(&user03Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&grid01Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&grid02Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&grid03Uuid)
	var row17Uuid, rowIntUuid, row23Uuid string
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", grid01Uuid, "test-17").Scan(&row17Uuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and int1= $2", grid02Uuid, 100).Scan(&rowIntUuid)
	db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", grid03Uuid, "test-23").Scan(&row23Uuid)

	t.Run("RootCanGetGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":true,"gridSpecialAccess":false`)
	})

	t.Run("User01CanGetGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user01", user01Uuid, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":false,"canEditGrid":false,"gridSpecialAccess":true`)
	})

	t.Run("RootCannotGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":false,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User01CanGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user01", user01Uuid, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":true,"gridSpecialAccess":false`)
	})

	t.Run("User01CanGetGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user01", user01Uuid, "/test/api/v1/"+grid02Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":true,"gridSpecialAccess":false`)
	})

	t.Run("User02CannotGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user02", user02Uuid, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":false,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User02CannotGetGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user02", user02Uuid, "/test/api/v1/"+grid02Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":false,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User02CannotGetGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user02", user02Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":false,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User01SetAccessForGrid02", func(t *testing.T) {
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship2","fromUuid":"` + grid02Uuid + `","toGridUuid":"` + model.UuidAccessLevel + `","uuid":"` + model.UuidAccessLevelReadAccess + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
	})

	t.Run("User02CanGetGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user02", user02Uuid, "/test/api/v1/"+grid02Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User01SetViewAccessForUser2Grid02", func(t *testing.T) {
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship4","fromUuid":"` + grid03Uuid + `","toGridUuid":"` + model.UuidUsers + `","uuid":"` + user02Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
	})

	t.Run("User02CanGetGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user02", user02Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User03CannotGetGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user03", user03Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":false,"canEditGrid":false,"gridSpecialAccess":false`)
	})

	t.Run("User01SetEditAccessForUser3Grid03", func(t *testing.T) {
		postStr := `{"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship5","fromUuid":"` + grid03Uuid + `","toGridUuid":"` + model.UuidUsers + `","uuid":"` + user03Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
	})

	t.Run("User03CanGetGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user03", user03Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":true,"gridSpecialAccess":false`)
	})

	t.Run("User01CanGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "user01", user01Uuid, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewGrid":true,"canEditGrid":true,"gridSpecialAccess":false`)
	})
}
