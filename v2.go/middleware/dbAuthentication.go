// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"database/sql"

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
			"SELECT uuid, firstName, lastName FROM users WHERE id = $1 AND password = crypt($2, password)",
			id,
			password).
			Scan(&uuid, &firstName, &lastName); err != nil {
			if err == sql.ErrNoRows {
				utils.Log("[%q] Invalid ID or password.", dbName)
			} else {
				utils.Log("[%q] Unknown error.", dbName)
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
