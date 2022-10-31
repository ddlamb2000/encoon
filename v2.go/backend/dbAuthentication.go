// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"fmt"

	"d.lambert.fr/encoon/utils"
)

func isDbAuthorized(dbName string, id string, password string) (string, string, string, error) {
	db := getDbByName(dbName)
	if db != nil {
		var uuid, firstName, lastName string
		sql := "SELECT uuid, text01, text02 FROM rows WHERE gridUuid = $1 AND uri = $2 AND text03 = crypt($3, text03)"
		if err := db.QueryRow(sql, utils.UuidUsers, id, password).Scan(&uuid, &firstName, &lastName); err != nil {
			return "", "", "", fmt.Errorf("[%q] Invalid ID or password: %v", dbName, err)
		}
		utils.Log("[%q] ID and password verified.", dbName)
		return uuid, firstName, lastName, nil
	}
	return "", "", "", fmt.Errorf("[%q] No database connection", dbName)
}
