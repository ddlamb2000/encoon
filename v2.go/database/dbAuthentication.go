// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type apiAuthResponse struct {
	uuid      string
	firstName string
	lastName  string
	err       error
}

func IsDbAuthorized(dbName string, user string, password string) (string, string, string, bool, error) {
	if db, err := GetDbByName(dbName); err == nil {
		var uuid, firstName, lastName string
		selectSql := " SELECT uuid, text1, text2 FROM rows "
		whereSql := " WHERE gridUuid = $1 AND text1 = $2 AND text4 = crypt($3, text4) "
		ctxChan := make(chan apiAuthResponse, 1)
		ctx, cancel := configuration.GetContextWithTimeOut(dbName)
		defer cancel()
		go func() {
			if err := TestSleep(ctx, dbName, db); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", utils.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
				return
			}
			if err := db.
				QueryRowContext(ctx, selectSql+whereSql, model.UuidUsers, user, password).
				Scan(&uuid, &firstName, &lastName); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", utils.LogAndReturnError("[%s] Invalid username or passphrase for %q: %v.", dbName, user, err)}
				return
			}
			utils.Log("[%s] ID and password verified for %q.", dbName, user)
			ctxChan <- apiAuthResponse{uuid, firstName, lastName, nil}
		}()
		select {
		case <-ctx.Done():
			return "", "", "", true, utils.LogAndReturnError("[%s] [%s] Authentication request has been cancelled: %v.", dbName, user, ctx.Err())
		case response := <-ctxChan:
			return response.uuid, response.firstName, response.lastName, false, response.err
		}
	}
	return "", "", "", false, utils.LogAndReturnError("[%s] No database connection.", dbName)
}
