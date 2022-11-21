// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

type apiAuthResponse struct {
	uuid      string
	firstName string
	lastName  string
	err       error
}

func IsDbAuthorized(ct context.Context, dbName string, user string, password string) (string, string, string, bool, error) {
	if db, err := GetDbByName(dbName); err == nil {
		var uuid, firstName, lastName string
		selectSql := " SELECT uuid, text1, text2 FROM rows "
		whereSql := " WHERE gridUuid = $1 AND text1 = $2 AND text4 = crypt($3, text4) "
		ctxChan := make(chan apiAuthResponse, 1)
		ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
		defer cancel()
		go func() {
			if err := TestSleep(ctx, dbName, db); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", configuration.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
				return
			}
			if err := db.
				QueryRowContext(ctx, selectSql+whereSql, model.UuidUsers, user, password).
				Scan(&uuid, &firstName, &lastName); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", configuration.LogAndReturnError("[%s] Invalid username or passphrase for %q: %v.", dbName, user, err)}
				return
			}
			configuration.Log("[%s] ID and password verified for %q.", dbName, user)
			ctxChan <- apiAuthResponse{uuid, firstName, lastName, nil}
		}()
		select {
		case <-ctx.Done():
			return "", "", "", true, configuration.LogAndReturnError("[%s] [%s] Authentication request has been cancelled: %v.", dbName, user, ctx.Err())
		case response := <-ctxChan:
			return response.uuid, response.firstName, response.lastName, false, response.err
		}
	}
	return "", "", "", false, configuration.LogAndReturnError("[%s] No database connection.", dbName)
}
