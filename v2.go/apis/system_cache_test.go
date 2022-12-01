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

func RunSystemTestCache(t *testing.T) {
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
}
