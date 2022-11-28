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
		request := "SELECT uuid, text1, text2 FROM users WHERE gridUuid = $1 AND enabled = true AND text1 = $2 AND text4 = crypt($3, text4)"
		ctxChan := make(chan apiAuthResponse, 1)
		ctx, cancel := configuration.GetContextWithTimeOut(ct, dbName)
		defer cancel()
		go func() {
			Sleep(ctx, dbName, user, db)
			if err := db.QueryRowContext(ctx, request, model.UuidUsers, user, password).Scan(&uuid, &firstName, &lastName); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", configuration.LogAndReturnError(dbName, user, "Invalid username or passphrase: %v.", err)}
				return
			}
			configuration.Log(dbName, user, "ID and password verified.")
			ctxChan <- apiAuthResponse{uuid, firstName, lastName, nil}
		}()
		select {
		case <-ctx.Done():
			return "", "", "", true, configuration.LogAndReturnError(dbName, user, "Authentication request has been cancelled: %v.", ctx.Err())
		case response := <-ctxChan:
			return response.uuid, response.firstName, response.lastName, false, response.err
		}
	}
	return "", "", "", false, configuration.LogAndReturnError(dbName, user, "No database connection.")
}
