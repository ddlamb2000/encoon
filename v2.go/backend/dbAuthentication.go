// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"d.lambert.fr/encoon/utils"
)

type apiAuthResponse struct {
	uuid      string
	firstName string
	lastName  string
	err       error
}

func isDbAuthorized(dbName string, user string, password string) (string, string, string, bool, error) {
	if db := getDbByName(dbName); db != nil {
		var uuid, firstName, lastName string
		selectSql := " SELECT uuid, text01, text02 FROM rows "
		whereSql := " WHERE gridUuid = $1 AND text01 = $2 AND text04 = crypt($3, text04) "
		ctxChan := make(chan apiAuthResponse, 1)
		ctx, cancel := utils.GetContextWithTimeOut()
		defer cancel()
		go func() {
			if err := db.
				QueryRowContext(ctx, selectSql+whereSql, utils.UuidUsers, user, password).
				Scan(&uuid, &firstName, &lastName); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", utils.LogAndReturnError("[%s] Invalid username or passphrase for %q: %v.", dbName, user, err)}
				return
			}
			if err := testSleep(ctx, dbName, db); err != nil {
				ctxChan <- apiAuthResponse{"", "", "", utils.LogAndReturnError("[%s] [%s] Sleep interrupted: %v.", dbName, user, err)}
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
