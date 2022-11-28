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
	var row17Uuid, rowInt100Uuid, row23Uuid string
	db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid01Uuid, "test-17").Scan(&row17Uuid)
	db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and int1= $2", grid02Uuid, 100).Scan(&rowInt100Uuid)
	db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid03Uuid, "test-23").Scan(&row23Uuid)

	t.Run("RootCanGetGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":true`)
	})

	t.Run("User01CanGetGrid", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canAddRows":true`)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditOwnedRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":false`)
	})

	t.Run("RootCannotGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "root", model.UuidRootUser, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringContains(t, responseData, `Access forbidden`)
	})

	t.Run("User01CanGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canAddRows":true`)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":true`)
	})

	t.Run("User01CanGetGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid02Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canAddRows":true`)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":true`)
	})

	t.Run("User02CannotGetGrid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid01Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringContains(t, responseData, `Access forbidden`)
	})

	t.Run("User02CannotGetGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid02Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringContains(t, responseData, `Access forbidden`)
	})

	t.Run("User02CannotGetGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringContains(t, responseData, `Access forbidden`)
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
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid02Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canAddRows":false`)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":false`)
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
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canAddRows":false`)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":false`)
	})

	t.Run("User03CannotGetGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringContains(t, responseData, `Access forbidden`)
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
		responseData, code, err := runGETRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid03Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"canAddRows":true`)
		jsonStringContains(t, responseData, `"canViewRows":true`)
		jsonStringContains(t, responseData, `"canEditOwnedRows":true`)
		jsonStringContains(t, responseData, `"canEditRows":false`)
	})

	t.Run("User01CanGetRow17Grid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid01Uuid+"/"+row17Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"text1":"test-17"`)
	})

	t.Run("User01CanUpdateRow17Grid01", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + row17Uuid + `","text1":"test-17 {2}","text2":"test-18 {2}","text3":"test-19 {2}","text4":"test-20 {2}"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-17 {2}","text2":"test-18 {2}","text3":"test-19 {2}","text4":"test-20 {2}"`)
	})

	t.Run("User01CanAddRowsGrid01", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-20","text2":"test-21","text3":"test-22","text4":"test-23"},` +
			`{"text1":"test-24","text2":"test-25","text3":"test-26","text4":"test-27"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-20","text2":"test-21","text3":"test-22","text4":"test-23"`)
		jsonStringContains(t, responseData, `"text1":"test-24","text2":"test-25","text3":"test-26","text4":"test-27"`)
	})

	t.Run("User01CanDeleteRowsGrid01", func(t *testing.T) {
		var row24Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid01Uuid, "test-24").Scan(&row24Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + row24Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringDoesntContain(t, responseData, `"text1":"test-24"`)
	})

	t.Run("User01CanGetRowGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid02Uuid+"/"+rowInt100Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"int1":100,"int2":100,"int3":100,"int4":100`)
	})

	t.Run("User01CanUpdateRowGrid02", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + rowInt100Uuid + `","int1":101,"int2":101,"int3":101,"int4":101}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"int1":101,"int2":101,"int3":101,"int4":101`)
	})

	t.Run("User01CanAddRowsGrid02", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"int1":200,"int2":200,"int3":200,"int4":200},` +
			`{"int1":300,"int2":300,"int3":300,"int4":300}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"int1":200,"int2":200,"int3":200,"int4":200`)
		jsonStringContains(t, responseData, `"int1":300,"int2":300,"int3":300,"int4":300`)
	})

	t.Run("User01CanDeleteRowsGrid02", func(t *testing.T) {
		var rowInt300Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and int1= $2", grid02Uuid, 300).Scan(&rowInt300Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + rowInt300Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringDoesntContain(t, responseData, `"int1":300`)
	})

	t.Run("User01CanGetRowGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid03Uuid+"/"+row23Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"text1":"test-23"`)
	})

	t.Run("User01CanUpdateRowGrid03", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + row23Uuid + `","text1":"test-23 {2}","text2":"test-24 {2}","text3":"test-25 {2}","text4":"test-26 {2}","int1":27,"int2":28,"int3":29,"int4":30}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-23 {2}","text2":"test-24 {2}","text3":"test-25 {2}","text4":"test-26 {2}","int1":27,"int2":28,"int3":29,"int4":30`)
	})

	t.Run("User01CanAddRowsGrid03", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a", "text1":"test-31","text2":"test-32","text3":"test-33","text4":"test-34","int1":35,"int2":36,"int3":37,"int4":38},` +
			`{"uuid":"b", "text1":"test-39","text2":"test-40","text3":"test-41","text4":"test-42","int1":43,"int2":44,"int3":45,"int4":46}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + grid01Uuid + `","uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-31","text2":"test-32","text3":"test-33","text4":"test-34","int1":35,"int2":36,"int3":37,"int4":38`)
		jsonStringContains(t, responseData, `"text1":"test-39","text2":"test-40","text3":"test-41","text4":"test-42","int1":43,"int2":44,"int3":45,"int4":46`)
	})

	t.Run("User01CanDeleteRowsGrid03", func(t *testing.T) {
		var row39Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid03Uuid, "test-39").Scan(&row39Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + row39Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringDoesntContain(t, responseData, `"text1":"test-39"`)
	})

	t.Run("User02CannotGetRow17Grid01", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid01Uuid+"/"+row17Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringDoesntContain(t, responseData, `"text1":"test-17"`)
	})

	t.Run("User02CannotUpdateRow17Grid01", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + row17Uuid + `","text1":"test-17 {3}","text2":"test-18 {3}","text3":"test-19 {2}","text4":"test-20 {2}"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringDoesntContain(t, responseData, `"text1":"test-17 {3}"`)
	})

	t.Run("User02CannotAddRowsGrid01", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-25","text2":"test-21","text3":"test-22","text4":"test-23"},` +
			`{"text1":"test-24","text2":"test-25","text3":"test-26","text4":"test-27"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringDoesntContain(t, responseData, `"text1":"test-25"`)
	})

	t.Run("User02CannotDeleteRowsGrid01", func(t *testing.T) {
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CanGetRowGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid02Uuid+"/"+rowInt100Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"int1":101,"int2":101,"int3":101,"int4":101`)
	})

	t.Run("User02CannotUpdateRowGrid02", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + rowInt100Uuid + `","int1":102,"int2":102,"int3":102,"int4":102}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CannotAddRowsGrid02", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"int1":400,"int2":200,"int3":200,"int4":200},` +
			`{"int1":500,"int2":300,"int3":300,"int4":300}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CannotDeleteRowsGrid02", func(t *testing.T) {
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + rowInt100Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CanGetRowGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid03Uuid+"/"+row23Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"text1":"test-23 {2}"`)
	})

	t.Run("User02CannotUpdateRowGrid03", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + row23Uuid + `","text1":"test-23 {4}","text2":"test-24 {4}","text3":"test-25 {3}","text4":"test-26 {3}","int1":27,"int2":28,"int3":29,"int4":30}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CannotAddRowsGrid03", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a", "text1":"test-40","text2":"test-32","text3":"test-33","text4":"test-34","int1":35,"int2":36,"int3":37,"int4":38},` +
			`{"uuid":"b", "text1":"test-41","text2":"test-40","text3":"test-41","text4":"test-42","int1":43,"int2":44,"int3":45,"int4":46}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + grid01Uuid + `","uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test01", user02Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CannotDeleteRowsGrid03", func(t *testing.T) {
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + row23Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CannotGetRow17Grid01", func(t *testing.T) {
		_, code, err := runGETRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid01Uuid+"/"+row17Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User02CannotUpdateRow17Grid01", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + row17Uuid + `","text1":"test-17 {6}","text2":"test-18 {2}","text3":"test-19 {2}","text4":"test-20 {2}"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CannotAddRowsGrid01", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-28","text2":"test-21","text3":"test-22","text4":"test-23"},` +
			`{"text1":"test-24","text2":"test-25","text3":"test-26","text4":"test-27"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CannotDeleteRowsGrid01", func(t *testing.T) {
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid01Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CanGetRowGrid02", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid02Uuid+"/"+rowInt100Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"int1":101,"int2":101,"int3":101,"int4":101`)
	})

	t.Run("User03CannotUpdateRowGrid02", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + rowInt100Uuid + `","int1":601,"int2":101,"int3":101,"int4":101}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CannotAddRowsGrid02", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"int1":800,"int2":200,"int3":200,"int4":200},` +
			`{"int1":300,"int2":300,"int3":300,"int4":300}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CannotDeleteRowsGrid02", func(t *testing.T) {
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + rowInt100Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid02Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CanGetRowGrid03", func(t *testing.T) {
		responseData, code, err := runGETRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid03Uuid+"/"+row23Uuid)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusOK)
		jsonStringContains(t, responseData, `"text1":"test-23 {2}"`)
	})

	t.Run("User03CannotUpdateRowGrid03", func(t *testing.T) {
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + row23Uuid + `","text1":"test-23 {7}","text2":"test-24 {7}","text3":"test-25 {7}","text4":"test-26 {7}","int1":27,"int2":28,"int3":29,"int4":30}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User03CanAddRowsGrid03", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a", "text1":"test-47","text2":"test-48","text3":"test-49","text4":"test-50","int1":51,"int2":52,"int3":53,"int4":54},` +
			`{"uuid":"b", "text1":"test-55","text2":"test-56","text3":"test-57","text4":"test-58","int1":59,"int2":60,"int3":61,"int4":62}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + grid01Uuid + `","uuid":"` + row17Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-47","text2":"test-48","text3":"test-49","text4":"test-50","int1":51,"int2":52,"int3":53,"int4":54`)
		jsonStringContains(t, responseData, `"text1":"test-55","text2":"test-56","text3":"test-57","text4":"test-58","int1":59,"int2":60,"int3":61,"int4":62`)
	})

	t.Run("User03CannotDeleteRowsGrid03", func(t *testing.T) {
		var row55Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", grid03Uuid, "test-55").Scan(&row55Uuid)
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + row55Uuid + `"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test03", user03Uuid, "/test/api/v1/"+grid03Uuid, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User01CannotCreateUser", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"testxx","text2":"Zero-xxx","text3":"Test xxx","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidUsers, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
		jsonStringContains(t, responseData, `Access forbidden`)
	})

	t.Run("User01CreateNewColumnsFor4thGrid", func(t *testing.T) {
		var gridUuid3 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)

		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Test Column 21","text2":"text1"},` +
			`{"uuid":"b","text1":"Test Column 22","text2":"text2"},` +
			`{"uuid":"c","text1":"Test Column 23","text2":"relationship1"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship1","fromUuid":"b","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidTextColumnType + `"},` +
			`{"columnName":"relationship2","fromUuid":"c","toGridUuid":"` + model.UuidGrids + `","uuid":"` + gridUuid3 + `"},` +
			`{"columnName":"relationship1","fromUuid":"c","toGridUuid":"` + model.UuidColumnTypes + `","uuid":"` + model.UuidReferenceColumnType + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidColumns, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Test Column 21","text2":"text1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 22","text2":"text2"`)
	})

	t.Run("User01Create4thSingleGrid", func(t *testing.T) {
		var column21Uuid, column22Uuid, column23Uuid string
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 21").Scan(&column21Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 22").Scan(&column22Uuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 23").Scan(&column23Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"Grid04","text2":"Test grid 04","text3":"journal"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column21Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column22Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + model.UuidColumns + `","uuid":"` + column23Uuid + `"},` +
			`{"columnName":"relationship2","fromUuid":"a","toGridUuid":"` + model.UuidAccessLevel + `","uuid":"` + model.UuidAccessLevelWriteAccess + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid04","text2":"Test grid 04","text3":"journal"`)
	})

	t.Run("User01CanModify4thSingleGrid", func(t *testing.T) {
		var gridUuid4 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid04").Scan(&gridUuid4)
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + gridUuid4 + `","text1":"Grid04","text2":"Test grid 04","text3":"person"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"Grid04","text2":"Test grid 04","text3":"person"`)
	})

	t.Run("User02CannotModify4thSingleGrid", func(t *testing.T) {
		var gridUuid4 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid04").Scan(&gridUuid4)
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + gridUuid4 + `","text1":"Grid04","text2":"Test grid 04","text3":"grid"}` +
			`]` +
			`}`
		_, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+model.UuidGrids, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusForbidden)
	})

	t.Run("User01CanAddRowIn4thSingleGrid", func(t *testing.T) {
		var gridUuid3, gridUuid4 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid04").Scan(&gridUuid4)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)
		var row09Uuid, row16Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", gridUuid3, "test-09").Scan(&row09Uuid)
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", gridUuid3, "test-16").Scan(&row16Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"test-01","text2":"test-02"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"},` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + gridUuid3 + `","uuid":"` + row16Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test01", user01Uuid, "/test/api/v1/"+gridUuid4, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-01","text2":"test-02"`)
	})

	t.Run("User02CanAddRowIn4thSingleGrid", func(t *testing.T) {
		var gridUuid3, gridUuid4 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid04").Scan(&gridUuid4)
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid03").Scan(&gridUuid3)
		var row09Uuid string
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", gridUuid3, "test-09").Scan(&row09Uuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"uuid":"a","text1":"test-03","text2":"test-04"}` +
			`],` +
			`"referencedValuesAdded":` +
			`[` +
			`{"columnName":"relationship1","fromUuid":"a","toGridUuid":"` + gridUuid3 + `","uuid":"` + row09Uuid + `"}` +
			`]` +
			`}`
		responseData, code, err := runPOSTRequestForUser("test", "test02", user02Uuid, "/test/api/v1/"+gridUuid4, postStr)
		errorIsNil(t, err)
		httpCodeEqual(t, code, http.StatusCreated)
		jsonStringContains(t, responseData, `"text1":"test-03","text2":"test-04"`)
	})
}
