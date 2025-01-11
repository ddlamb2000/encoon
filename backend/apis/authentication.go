// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
)

func handleAuthentication(dbName string, content requestContent) responseContent {
	if dbName == "" || content.Userid == "" || content.Password == "" {
		return responseContent{
			Status:      FailedStatus,
			Action:      content.Action,
			TextMessage: "Authentication: missing username or passphrase",
		}
	}
	configuration.Log(dbName, "", "Authentication: %s try to login", content.Userid)
	userUuid, firstName, lastName, timeOut, err := database.IsDbAuthorized(context.Background(), dbName, content.Userid, content.Password)
	if err != nil || userUuid == "" {
		if timeOut {
			configuration.LogError(dbName, "", "Authentication: time out ", err)
			return responseContent{
				Status:      FailedStatus,
				Action:      content.Action,
				ActionText:  content.ActionText,
				TextMessage: "Authentication: time out " + err.Error(),
			}
		} else {
			configuration.LogError(dbName, "", "Authentication: failed ", err)
			return responseContent{
				Status:      FailedStatus,
				Action:      content.Action,
				ActionText:  content.ActionText,
				TextMessage: "Authentication: failed " + err.Error(),
			}
		}
	}
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
	token, err := getNewToken(dbName, content.Userid, userUuid, firstName, lastName, expiration)
	if err != nil {
		configuration.LogError(dbName, "", "Authentication: creation of JWT failed ", err)
		return responseContent{
			Status:      FailedStatus,
			Action:      content.Action,
			ActionText:  content.ActionText,
			TextMessage: "Authentication: creation of JWT failed " + err.Error(),
		}
	}
	configuration.Log(dbName, content.Userid, "Connected.")
	return responseContent{
		Status:      SuccessStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		FirstName:   firstName,
		LastName:    lastName,
		TextMessage: "User authenticated",
		JWT:         token,
	}
}

// function is available for mocking
var getNewToken = func(dbName, user, userUuid, firstName, lastName string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user":          user,
		"userUuid":      userUuid,
		"userFirstName": firstName,
		"userLastName":  lastName,
		"expires":       expiration,
	})
	configuration.Trace(dbName, user, "Token generated, expiration: %v", expiration)
	jwtSecret := configuration.GetJWTSecret(dbName)
	return token.SignedString([]byte(jwtSecret))
}

func getTokenParsingHandler(dbName string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if ok := verifyToken(token); !ok {
			return nil, configuration.LogAndReturnError(dbName, "", "Unexpect signing method: %v.", token.Header["alg"])
		}
		return []byte(configuration.GetJWTSecret(dbName)), nil
	}
}

// function is available for mocking
var verifyToken = func(token *jwt.Token) bool {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	return ok
}
