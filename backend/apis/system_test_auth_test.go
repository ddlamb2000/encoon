// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"errors"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
	"github.com/golang-jwt/jwt/v5"
)

func RunSystemTestAuth(t *testing.T) {
	t.Run("AuthInvalid1", func(t *testing.T) {
		response, responseData := runKafkaTestAuthRequest(t, "test", ApiParameters{
			Action: ActionAuthentication,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Missing username or passphrase"`)
	})

	t.Run("AuthInvalid2", func(t *testing.T) {
		response, responseData := runKafkaTestAuthRequest(t, "test", ApiParameters{
			Action:   ActionAuthentication,
			Userid:   "root",
			Password: "======",
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Authentication failed"`)
	})

	t.Run("AuthValid", func(t *testing.T) {
		response, responseData := runKafkaTestAuthRequest(t, "test", ApiParameters{
			Action:   ActionAuthentication,
			Userid:   "root",
			Password: "dGVzdA==",
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"textMessage":"User authenticated","firstName":"root","lastName":"root","jwt":`)
	})

	t.Run("ApiUsersWithTimeOut", func(t *testing.T) {
		database.ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
		defer setDefaultTestSleepTimeAndTimeOutThreshold()
		response, responseData := runKafkaTestAuthRequest(t, "test", ApiParameters{
			Action:   ActionAuthentication,
			Userid:   "root",
			Password: "dGVzdA==",
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Authentication timed out"`)
	})

	t.Run("ApiUsersNoHeader", func(t *testing.T) {
		response, responseData := runKafkaTestAuthRequest(t, "test", ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
			Uuid:     model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Missing authorization"`)
	})

	t.Run("ApiUsersIncorrectToken", func(t *testing.T) {
		response, responseData := runKafkaTestRequestWithToken(t, "test", "root", model.UuidRootUser, model.UuidUsers, "xxxxxxxxxxx", ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidUsers,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Missing authorization"`)
	})

	t.Run("ApiUsersIncorrectToken2", func(t *testing.T) {
		getNewTokenImpl := getNewToken
		getNewToken = func(dbName, user, userUuid, firstName, lastName string, expiration time.Time) (string, error) {
			return "", errors.New("xxx")
		} // mock function
		response, responseData := runKafkaTestAuthRequest(t, "test", ApiParameters{
			Action:   ActionAuthentication,
			Userid:   "root",
			Password: "dGVzdA==",
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Creation of JWT failed"`)
		getNewToken = getNewTokenImpl
	})

	t.Run("ApiUsersExpired", func(t *testing.T) {
		expiration := time.Now().Add(-time.Duration(configuration.GetConfiguration().JwtExpiration) * time.Minute)
		token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
		response, responseData := runKafkaTestRequestWithToken(t, "test", "root", model.UuidRootUser, model.UuidUsers, token, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidUsers,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Authorization expired"`)
	})

	t.Run("ApiUsersPassing", func(t *testing.T) {
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().JwtExpiration) * time.Minute)
		token, _ := getNewToken("test", "root", model.UuidRootUser, "root", "root", expiration)
		response, responseData := runKafkaTestRequestWithToken(t, "test", "root", model.UuidRootUser, model.UuidUsers, token, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"text1":"root","text2":"root"`)
	})

	t.Run("ApiUsersIncorrectToken3", func(t *testing.T) {
		getNewTokenImpl := getNewToken
		getNewToken = func(dbName, user, userUuid, firstName, lastName string, expiration time.Time) (string, error) {
			return "", errors.New("xxx")
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Missing authorization"`)
		getNewToken = getNewTokenImpl
	})

	t.Run("ApiUsersNotFound2", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "d7c004ff-cccc-dddd-eeee-cd42b2847508", ApiParameters{
			Action:   ActionLoad,
			GridUuid: "d7c004ff-cccc-dddd-eeee-cd42b2847508",
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Data not found"`)
	})

	t.Run("Post404", func(t *testing.T) {
		postStr := `{}`
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "", ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: "",
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `Error when retrieving grid definition: pq: invalid input syntax for type uuid`)
	})

	t.Run("CreateUserNoData", func(t *testing.T) {
		postStr := `{}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", model.UuidRootUser, model.UuidUsers, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidUsers,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidUsers+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `{"gridUuid":"`+model.UuidUsers+`","uuid":"`+model.UuidRootUser+`"`)
	})

	t.Run("CreateNewSingleUser", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test01","text2":"Zero-one","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", model.UuidRootUser, model.UuidUsers, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidUsers,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringDoesntContain(t, responseData, `"countRows":1`)
		jsonStringContains(t, responseData, `"countRows":2`)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidUsers+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `{"gridUuid":"`+model.UuidUsers+`","uuid":"`+model.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"text1":"test01","text2":"Zero-one","text3":"Test"`)
	})

	t.Run("Create3Users", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test02","text2":"Zero-two","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"},` +
			`{"text1":"test03","text2":"Zero-three","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"},` +
			`{"text1":"test04","text2":"Zero-four","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		response, responseData := runKafkaTestRequest(t, "test", "test01", model.UuidRootUser, model.UuidUsers, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidUsers,
			DataSet:  stringToJson(postStr),
		})
		responseIsSuccess(t, response)
		jsonStringDoesntContain(t, responseData, `"countRows":2`)
		jsonStringContains(t, responseData, `"countRows":5`)
		jsonStringContains(t, responseData, `"text1":"test02","text2":"Zero-two","text3":"Test"`)
		jsonStringContains(t, responseData, `"text1":"test03","text2":"Zero-three","text3":"Test"`)
		jsonStringContains(t, responseData, `"text1":"test04","text2":"Zero-four","text3":"Test"`)
	})

	t.Run("CreateWithIncorrectUserUuid", func(t *testing.T) {
		postStr := `{"rowsAdded":` +
			`[` +
			`{"text1":"test02","text2":"Zero-two","text3":"Test","text4":"$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"}` +
			`]` +
			`}`
		response, _ := runKafkaTestRequest(t, "test", "root", "xxyyzz", model.UuidUsers, ApiParameters{
			Action:   ActionChangeGrid,
			GridUuid: model.UuidUsers,
			DataSet:  stringToJson(postStr),
		})
		responseIsFailure(t, response)
	})

	t.Run("ApiUsersDefectToken", func(t *testing.T) {
		verifyTokenImpl := verifyToken
		verifyToken = func(*jwt.Token) bool { return false } // mock function
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().JwtExpiration) * time.Minute)
		token, _ := getNewToken("test", "root", model.UuidRootUser, "root", "root", expiration)
		response, responseData := runKafkaTestRequestWithToken(t, "test", "root", model.UuidRootUser, model.UuidUsers, token, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidUsers,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Missing authorization"`)
		verifyToken = verifyTokenImpl
	})
}
