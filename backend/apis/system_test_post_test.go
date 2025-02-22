// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"errors"
	"testing"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func RunSystemTestPost(t *testing.T) {
	var user01Uuid string
	db, _ := database.GetDbByName("test")
	db.QueryRow("SELECT uuid FROM users WHERE gridUuid = $1 and text1= $2", model.UuidUsers, "test01").Scan(&user01Uuid)

	t.Run("VerifyPostDbNotConfigured", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "baddb", "root", model.UuidRootUser, "xxx", ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: "xxx",
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Unable to connect to database: dial tcp`)
	})

	t.Run("CreateNewSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"journal"`)
	})

	t.Run("VerifySingleGridIsCreated", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"journal"`)
	})

	t.Run("VerifyNoRowInSingleGrid", func(t *testing.T) {
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionLoad,
			GridUuid: gridUuid,
		})
		responseIsSuccess(t, response)
		jsonStringDoesntContain(t, responseData, `"countRows":`)
		jsonStringDoesntContain(t, responseData, `"rows"`)
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
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidColumns,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Test Column 01","text2":"text1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 02","text2":"text2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 03","text2":"text3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 04","text2":"text4"`)

		var gridUuid, uuidCol1, uuidCol2, uuidCol3, uuidCol4 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 01").Scan(&uuidCol1)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 02").Scan(&uuidCol2)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 03").Scan(&uuidCol3)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 04").Scan(&uuidCol4)

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
		response, _ = runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidRelationships, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidRelationships,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
	})

	t.Run("CreateNewRowInSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test01","text2":"test02","text3":"test03","text4":"test04"}` +
			`]` +
			`}`
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Test Column 01","name":"text1","type":"Text"`)
		jsonStringContains(t, responseData, `"label":"Test Column 02","name":"text2","type":"Text"`)
	})

	t.Run("CreateRowIncorrectPayload", func(t *testing.T) {
		postStr := `{"xxxxx"}`
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
	})

	t.Run("VerifyNewRowInSingleGrid", func(t *testing.T) {
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "root", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionLoad,
			GridUuid: gridUuid,
		})
		responseIsSuccess(t, response)
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
		defer setDefaultTestSleepTimeAndTimeOutThreshold()
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Post request has been cancelled: context deadline exceeded."`)
	})

	t.Run("VerifyNoNewRowInSingleGrid", func(t *testing.T) {
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "root", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionLoad,
			GridUuid: gridUuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"test02","text3":"test03","text4":"test04"`)
		jsonStringDoesntContain(t, responseData, `"text1":"test05"`)
	})

	t.Run("UpdateWrongRow", func(t *testing.T) {
		var uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + model.UuidUserColumnId + `","text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error retrieving row`)
		var revision int
		db.QueryRow("SELECT revision FROM grids WHERE gridUuid = $1 and uuid = $2", model.UuidGrids, uuid).Scan(&revision)
		intEqual(t, revision, 1)
	})

	t.Run("UpdateNewRow", func(t *testing.T) {
		var uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuid + `","text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"uuid":"`+uuid+`"`)
		jsonStringContains(t, responseData, `"text1":"Grid01","text2":"Test grid 01","text3":"journal"`)
		var revision int
		db.QueryRow("SELECT revision FROM grids WHERE gridUuid = $1 and uuid = $2", model.UuidGrids, uuid).Scan(&revision)
		intEqual(t, revision, 2)
	})

	t.Run("CreateNewandUpdateRowsInSingleGrid", func(t *testing.T) {
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
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
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringContains(t, responseData, `"text1":"test-01","text2":"test-02","text3":"test-03","text4":"test-04"`)
		jsonStringContains(t, responseData, `"text1":"test-09","text2":"test-10","text3":"test-11","text4":"test-12"`)
	})

	t.Run("DeleteWrongRowInSingleGrid", func(t *testing.T) {
		var uuidGrid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + model.UuidGridColumnName + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error retrieving row`)
	})

	t.Run("CreateDeleteRowsInSingleGrid", func(t *testing.T) {
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
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
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text1":"test-21","text2":"test-22","text3":"test-23","text4":"test-24"`)
		jsonStringContains(t, responseData, `"text1":"test-25","text2":"test-26","text3":"test-27","text4":"test-28"`)
		jsonStringDoesntContain(t, responseData, `"text1":"test-13"`)
	})

	t.Run("DeleteRowInSingleGrid", func(t *testing.T) {
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-09").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringDoesntContain(t, responseData, `"text1":"test-09"`)
	})

	t.Run("CreateNewRowInSingleGridWithTimeOut2", func(t *testing.T) {
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test-29","text2":"test-30","text3":"test-31","text4":"test-32"}` +
			`]` +
			`}`
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 10, 500)
		defer setDefaultTestSleepTimeAndTimeOutThreshold()
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":6`)
		jsonStringContains(t, responseData, `"text1":"test-29","text2":"test-30","text3":"test-31","text4":"test-32"`)
	})

	t.Run("Create2ndSingleGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid02","text2":"Test grid 02","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
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
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidColumns, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidColumns,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"Test Column 05","text2":"int1"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 06","text2":"int2"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 07","text2":"int3"`)
		jsonStringContains(t, responseData, `"text1":"Test Column 08","text2":"int4"`)

		var gridUuid, uuidCol1, uuidCol2, uuidCol3, uuidCol4 string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 05").Scan(&uuidCol1)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 06").Scan(&uuidCol2)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 07").Scan(&uuidCol3)
		db.QueryRow("SELECT uuid FROM columns WHERE gridUuid = $1 and text1= $2", model.UuidColumns, "Test Column 08").Scan(&uuidCol4)

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
		response, _ = runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidRelationships, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidRelationships,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
	})

	t.Run("CreateNewRowsIn2ndSingleGrid", func(t *testing.T) {
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid02").Scan(&gridUuid)
		postStr := `{"rowsAdded":` +
			`[` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":100,"int2":100,"int3":100,"int4":100},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4},` +
			`{"int1":1,"int2":2,"int3":3,"int4":4}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, gridUuid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: gridUuid,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":12`)
		jsonStringContains(t, responseData, `"int1":1,"int2":2,"int3":3,"int4":4`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Test Column 05","name":"int1","type":"Integer"`)
	})

	t.Run("InvalidCreateGrid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", model.UuidRootUser, "d7c004ff-cccc-dddd-eeee-cd42b2847508", ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: "d7c004ff-cccc-dddd-eeee-cd42b2847508",
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Data not found."`)
	})

	t.Run("InvalidCreateGrid2", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", model.UuidRootUser, "xxx", ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: "xxx",
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when retrieving grid definition: pq: invalid input syntax for type uuid`)
	})

	t.Run("CreateSingleGridDefect", func(t *testing.T) {
		getBeginTransactionQueryImpl := getBeginTransactionQuery
		getBeginTransactionQuery = func() string { return "xxx" } // mock function
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid0x","text2":"Test grid 0x","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Begin transaction error: pq: syntax error`)
		getBeginTransactionQuery = getBeginTransactionQueryImpl
	})

	t.Run("CreateSingleGridDefect2", func(t *testing.T) {
		getCommitTransactionQueryImpl := getCommitTransactionQuery
		getCommitTransactionQuery = func() string { return "xxx" } // mock function
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid0x","text2":"Test grid 0x","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Commit transaction error: pq: syntax error `)
		getCommitTransactionQuery = getCommitTransactionQueryImpl
	})

	t.Run("CreateSingleGridDefect3", func(t *testing.T) {
		getInsertStatementForReferenceRowImpl := getInsertStatementForReferenceRow
		getInsertStatementForReferenceRow = func() string { return "xxx" } // mock function
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid0x","text2":"Test grid 0x","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Insert referenced row error: pq: syntax error`)
		getInsertStatementForReferenceRow = getInsertStatementForReferenceRowImpl
	})

	t.Run("CreateSingleGridDefect4", func(t *testing.T) {
		getInsertStatementForGridsApiImpl := getInsertStatementForGridsApi
		getInsertStatementForGridsApi = func(grid *model.Grid) string { return "xxx" } // mock function
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"Grid0x","text2":"Test grid 0x","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Insert row error: pq: syntax error`)
		getInsertStatementForGridsApi = getInsertStatementForGridsApiImpl
	})

	t.Run("CreateSingleGridDefect5", func(t *testing.T) {
		getUpdateStatementForGridsApiImpl := getUpdateStatementForGridsApi
		getUpdateStatementForGridsApi = func(grid *model.Grid) string { return "xxx" } // mock function
		var uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuid + `","text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Update row error: pq: syntax error`)
		getUpdateStatementForGridsApi = getUpdateStatementForGridsApiImpl
	})

	t.Run("DeleteSingleGridDefect6", func(t *testing.T) {
		getDeleteGridReferencedRowQueryImpl := getDeleteGridReferencedRowQuery
		getDeleteGridReferencedRowQuery = func(*model.Grid) string { return "xxx" } // mock function
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-01").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Delete referenced row error: pq: syntax error`)
		getDeleteGridReferencedRowQuery = getDeleteGridReferencedRowQueryImpl
	})

	t.Run("DeleteSingleGridDefect7", func(t *testing.T) {
		getDeleteGridRowQueryImpl := getDeleteGridRowQuery
		getDeleteGridRowQuery = func(*model.Grid) string { return "xxx" } // mock function
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-01").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Delete row error: pq: syntax error`)
		getDeleteGridRowQuery = getDeleteGridRowQueryImpl
	})

	t.Run("DeleteSingleGridRowDefect", func(t *testing.T) {
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-29").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			Uuid:     "xxxx",
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Data not found."`)
	})

	t.Run("DeleteSingleGridRow", func(t *testing.T) {
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-29").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		postStr := `{"rowsDeleted":` +
			`[` +
			`{"uuid":"` + uuidRow + `"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: uuidGrid,
			Uuid:     uuidRow,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"enabled":false`)
	})

	t.Run("GetSingleGridDefect9", func(t *testing.T) {
		getGridForOwnershipImpl := getGridForOwnership
		getGridForOwnership = func(ApiRequest, *model.Grid, *model.Row) (*model.Grid, error) {
			return nil, errors.New("xxx")
		}
		var uuidGrid, uuidRow string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuidGrid)
		stringNotEqual(t, uuidGrid, "")
		db.QueryRow("SELECT uuid FROM rows WHERE gridUuid = $1 and text1= $2", uuidGrid, "test-29").Scan(&uuidRow)
		stringNotEqual(t, uuidRow, "")
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, uuidGrid, ApiParameters{
			Action:   ActionLoad,
			GridUuid: uuidGrid,
			Uuid:     uuidRow,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"xxx"`)
		getGridForOwnership = getGridForOwnershipImpl
	})

	t.Run("GetSingleGridDefect10", func(t *testing.T) {
		getRowsQueryForGridUuidAttachedToColumnImpl := getRowsQueryForGridUuidAttachedToColumn
		getRowsQueryForGridUuidAttachedToColumn = func() string { return "x x x" }
		response, responseData := runKafkaTestRequest(t, "test", "root", user01Uuid, model.UuidColumns, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidColumns,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when retrieving grid uuid for column`)
		getRowsQueryForGridUuidAttachedToColumn = getRowsQueryForGridUuidAttachedToColumnImpl
	})

	t.Run("VerifyAuditInSingleGrid", func(t *testing.T) {
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     gridUuid,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"audits":[`)
		jsonStringContains(t, responseData, `"createdByName":"test01","actionName":"Updated"`)
		jsonStringContains(t, responseData, `"createdByName":"test01","actionName":"Created"`)
	})

	t.Run("VerifyAuditInSingleGridDefect1", func(t *testing.T) {
		getAuditsQueryForRowImpl := getAuditsQueryForRow
		getAuditsQueryForRow = func(*model.Grid, string) string { return "x x x" }
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     gridUuid,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when querying audits: pq: syntax error`)
		getAuditsQueryForRow = getAuditsQueryForRowImpl
	})

	t.Run("VerifyAuditInSingleGridDefect2", func(t *testing.T) {
		getAuditsQueryOutputForRowImpl := getAuditsQueryOutputForRow
		getAuditsQueryOutputForRow = func(audit *model.Audit) []any { return nil }
		var gridUuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&gridUuid)
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     gridUuid,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when scanning audits: sql`)
		getAuditsQueryOutputForRow = getAuditsQueryOutputForRowImpl
	})

	t.Run("UpdateGridTransactionDefect", func(t *testing.T) {
		getInsertStatementForTransactionImpl := getInsertStatementForTransaction
		getInsertStatementForTransaction = func() string { return "xxx" } // mock function
		var uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuid + `","text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Insert transaction error: pq: syntax error at or near \"xxx\".`)
		getInsertStatementForTransaction = getInsertStatementForTransactionImpl
	})

	t.Run("UpdateGridTransactionDefect2", func(t *testing.T) {
		getInsertStatementForTransactionReferenceRowImpl := getInsertStatementForTransactionReferenceRow
		getInsertStatementForTransactionReferenceRow = func() string { return "xxx" } // mock function
		var uuid string
		db.QueryRow("SELECT uuid FROM grids WHERE gridUuid = $1 and text1= $2", model.UuidGrids, "Grid01").Scan(&uuid)
		stringNotEqual(t, uuid, "")
		postStr := `{"rowsEdited":` +
			`[` +
			`{"uuid":"` + uuid + `","text1":"Grid01","text2":"Test grid 01","text3":"journal"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", user01Uuid, model.UuidGrids, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidGrids,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Insert transaction referenced row error: pq: syntax error at or near \"xxx\"."`)
		getInsertStatementForTransactionReferenceRow = getInsertStatementForTransactionReferenceRowImpl
	})
}
