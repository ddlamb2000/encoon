// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"database/sql"

	"d.lambert.fr/encoon/backend"
	"d.lambert.fr/encoon/utils"
)

func isDbAuthorized(dbName string, id string, password string) (bool, string, string, string) {
	utils.Log("[%q] Verify ID and password.", dbName)
	db := dbs[dbName]
	if db != nil {
		var uuid string
		var firstName string
		var lastName string
		if err := db.QueryRow(
			"SELECT uuid, text02, text03 FROM rows WHERE gridUuid = $1 AND text01 = $2 AND text04 = crypt($3, text04)",
			backend.UuidUsers,
			id,
			password).
			Scan(&uuid, &firstName, &lastName); err != nil {
			if err == sql.ErrNoRows {
				utils.Log("[%q] Invalid ID or password: %v.", dbName, err)
			} else {
				utils.Log("[%q] Unknown error: %v.", dbName, err)
			}
		} else {
			if uuid != "" {
				utils.Log("[%q] ID and password verified.", dbName)
				return true, uuid, firstName, lastName
			} else {
				utils.Log("[%q] Invalid ID or password.", dbName)
			}
			utils.Log("[%q] Error.", err)
		}
	} else {
		utils.Log("[%q] No database connection.", dbName)
	}
	return false, "", "", ""
}
