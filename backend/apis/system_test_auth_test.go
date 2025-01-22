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

	t.Run("AuthInvalid3", func(t *testing.T) {
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// req, _ := http.NewRequest("POST", "/undefined/api/v1/authentication", strings.NewReader(`{"id": "root", "password": "======"}`))
		// w := httptest.NewRecorder()
		// testRouter.ServeHTTP(w, req)
		// responseData, err := io.ReadAll(w.Body)
		// httpCodeEqual(t, w.Code, http.StatusBadRequest)

		// expect := utils.CleanupStrings(`{"error":"No database parameter"}`)
		// response := utils.CleanupStrings(string(responseData))

		// if err != nil {
		// 	t.Errorf(`Response %v for %v: %v.`, response, w, err)
		// }
		// if response != expect {
		// 	t.Errorf(`Response %v incorrect instead of %v.`, response, expect)
		// }
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

	t.Run("404Html", func(t *testing.T) {
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// req, _ := http.NewRequest("GET", "/xxx/yyy", nil)
		// w := httptest.NewRecorder()
		// testRouter.ServeHTTP(w, req)
		// responseData, err := io.ReadAll(w.Body)
		// httpCodeEqual(t, w.Code, http.StatusNotFound)

		// expect := utils.CleanupStrings(`404 page not found`)
		// response := utils.CleanupStrings(string(responseData))

		// if err != nil {
		// 	t.Errorf(`Response %v for %v: %v.`, response, w, err)
		// }
		// if response != expect {
		// 	t.Errorf(`Response %v incorrect.`, response)
		// }
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
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// req, _ := http.NewRequest("GET", "/test/api/v1/"+model.UuidUsers, nil)
		// req.Header.Add("Authorization", "xxxxxxx")
		// w := httptest.NewRecorder()
		// testRouter.ServeHTTP(w, req)
		// responseData, err := io.ReadAll(w.Body)
		// httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		// expect := utils.CleanupStrings(`{"error":"Incorrect header."}`)
		// response := utils.CleanupStrings(string(responseData))

		// if err != nil {
		// 	t.Errorf(`Response %v for %v: %v.`, response, w, err)
		// }
		// if response != expect {
		// 	t.Errorf(`Response %v incorrect.`, response)
		// }
	})

	t.Run("ApiUsersMissingBearer", func(t *testing.T) {
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		// token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
		// req, _ := http.NewRequest("GET", "/test/api/v1/"+model.UuidUsers, nil)
		// req.Header.Add("Authorization", token)
		// w := httptest.NewRecorder()
		// testRouter.ServeHTTP(w, req)
		// responseData, err := io.ReadAll(w.Body)
		// httpCodeEqual(t, w.Code, http.StatusUnauthorized)

		// expect := utils.CleanupStrings(`Invalid request`)
		// response := utils.CleanupStrings(string(responseData))

		// if err != nil {
		// 	t.Errorf(`Response %v for %v: %v.`, response, w, err)
		// }
		// if !strings.Contains(response, expect) {
		// 	t.Errorf(`Response %v incorrect.`, response)
		// }
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

	t.Run("ApiUsersNotFound", func(t *testing.T) {
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		// token, _ := getNewToken("test", "root", "0", "root", "root", expiration)
		// req, _ := http.NewRequest("GET", "/test/api/v0/users", nil)
		// req.Header.Add("Authorization", "Bearer "+token)
		// w := httptest.NewRecorder()
		// testRouter.ServeHTTP(w, req)
		// responseData, err := io.ReadAll(w.Body)
		// httpCodeEqual(t, w.Code, http.StatusNotFound)

		// expect := utils.CleanupStrings(`404 page not found`)
		// response := utils.CleanupStrings(string(responseData))

		// if err != nil {
		// 	t.Errorf(`Response %v for %v: %v.`, response, w, err)
		// }
		// if !strings.Contains(response, expect) {
		// 	t.Errorf(`Response %v incorrect.`, response)
		// }
	})

	t.Run("ApiUsersNotFound2", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "d7c004ff-cccc-dddd-eeee-cd42b2847508", ApiParameters{
			Action:   ActionLoad,
			GridUuid: "d7c004ff-cccc-dddd-eeee-cd42b2847508",
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Data not found"`)
	})

	t.Run("ApiUsersNotFound3", func(t *testing.T) {
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		// token, _ := getNewToken("test", "root", model.UuidRootUser, "root", "root", expiration)
		// req, _ := http.NewRequest("GET", "/test/api/v1/us", nil)
		// req.Header.Add("Authorization", "Bearer "+token)
		// w := httptest.NewRecorder()
		// testRouter.ServeHTTP(w, req)
		// responseData, err := io.ReadAll(w.Body)
		// httpCodeEqual(t, w.Code, http.StatusInternalServerError)

		// expect := utils.CleanupStrings(`Error when retrieving grid definition: pq: invalid input syntax for type uuid`)
		// response := utils.CleanupStrings(string(responseData))

		// if err != nil {
		// 	t.Errorf(`Response %v for %v: %v.`, response, w, err)
		// }
		// if !strings.Contains(response, expect) {
		// 	t.Errorf(`Response %v incorrect instead of %v.`, response, expect)
		// }
	})

	t.Run("CreateNewUserNoAuth", func(t *testing.T) {
		// NOTE: MAYBE NOT APPLICABLE for Kafka

		// postStr := `{"rowsAdded":` +
		// 	`[` +
		// 	`{"text1":"aaaa","text2":"bbbb"}` +
		// 	`]` +
		// 	`}`
		// responseData, code, err := runPOSTRequestForUser("te st", "root", model.UuidRootUser, "/test/api/v1/"+model.UuidUsers, postStr)
		// errorIsNil(t, err)
		// httpCodeEqual(t, code, http.StatusUnauthorized)
		// jsonStringContains(t, responseData, `"error":"Invalid request or unauthorized database access: signature is invalid."`)
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
