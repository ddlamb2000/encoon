// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"time"

	"d.lambert.fr/encoon/utils"
)

func isDbAuthorized(dbName string, id string, password string) (string, string, string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(utils.Configuration.TimeOutThreshold)*time.Millisecond)
	defer cancel()
	if db := getDbByName(dbName); db != nil {
		var uuid, firstName, lastName string
		selectSql := " SELECT uuid, text01, text02 FROM rows "
		whereSql := " WHERE gridUuid = $1 AND text01 = $2 AND text04 = crypt($3, text04) "
		if err := db.
			QueryRowContext(ctx, selectSql+whereSql, utils.UuidUsers, id, password).
			Scan(&uuid, &firstName, &lastName); err != nil {
			return "", "", "", utils.LogAndReturnError("[%s] Invalid username or passphrase for %q: %v", dbName, id, err)
		}
		utils.Log("[%s] ID and password verified for %q.", dbName, id)
		return uuid, firstName, lastName, nil
	}
	return "", "", "", utils.LogAndReturnError("[%s] No database connection", dbName)
}
