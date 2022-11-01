// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"fmt"
	"time"

	"d.lambert.fr/encoon/utils"
)

func isDbAuthorized(dbName string, id string, password string) (string, string, string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(utils.Configuration.DbTimeOut)*time.Second)
	defer cancel()
	db := getDbByName(dbName)
	if db != nil {
		var uuid, firstName, lastName string
		selectSql := " SELECT uuid, text01, text02 "
		fromSql := " FROM rows "
		whereSql := " WHERE gridUuid = $1 AND uri = $2 AND text03 = crypt($3, text03) "
		if err := db.
			QueryRowContext(
				ctx,
				selectSql+fromSql+whereSql,
				utils.UuidUsers,
				id,
				password).
			Scan(
				&uuid,
				&firstName,
				&lastName); err != nil {
			return "", "", "", fmt.Errorf("[%q] Invalid ID or password: %v", dbName, err)
		}
		utils.Log("[%q] ID and password verified.", dbName)
		return uuid, firstName, lastName, nil
	}
	return "", "", "", fmt.Errorf("[%q] No database connection", dbName)
}
