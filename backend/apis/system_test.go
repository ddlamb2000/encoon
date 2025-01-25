// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

//
// Note
// System test can be run using command line `go test apis/* -v -run TestSystem`
//

package apis

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
)

func TestSystem(t *testing.T) {
	setDefaultTestSleepTimeAndTimeOutThreshold()
	InitializeCaches()
	t.Run("IncorrectDb", func(t *testing.T) { RunTestConnectDbServersIncorrect(t) })
	t.Run("ConnectDb", func(t *testing.T) { RunTestConnectDbServers(t) })
	t.Run("RecreateDb", func(t *testing.T) { RunTestRecreateDb(t) })

	configuration.LoadConfiguration("../testData/systemTest.yml")
	go ReadMessagesFromKafka()
	kafkaTestProducer = createKafkaTestProducer("test")
	kafkaTestConsumer = createKafkaTestConsumer("test")
	kafkaBaddbProducer = createKafkaTestProducer("baddb")
	kafkaBaddbConsumer = createKafkaTestConsumer("baddb")
	defer kafkaTestProducer.Close()
	defer kafkaTestConsumer.Close()
	defer kafkaBaddbProducer.Close()
	defer kafkaBaddbConsumer.Close()
	defer ShutdownKafkaProducers()
	defer ShutdownKafkaConsumers()
	time.Sleep(time.Second * 5) // wait for Kafka election

	t.Run("Kafka", func(t *testing.T) { RunSystemTestKafka(t) })
	t.Run("Auth", func(t *testing.T) { RunSystemTestAuth(t) })
	t.Run("Get", func(t *testing.T) { RunSystemTestGet(t) })
	t.Run("Post", func(t *testing.T) { RunSystemTestPost(t) })

	InitializeCaches()
	t.Run("PostRelationships", func(t *testing.T) { RunSystemTestPostRelationships(t) })
	t.Run("RowLevelAccess", func(t *testing.T) { RunSystemTestGetRowLevel(t) })
	t.Run("Cache", func(t *testing.T) { RunSystemTestCache(t) })
	t.Run("NotOwnedColumn", func(t *testing.T) { RunSystemTestNotOwnedColumn(t) })
}

func getTokenForUser(dbName, userName, userUuid string) string {
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().JwtExpiration) * time.Minute)
	token, _ := getNewToken(dbName, userName, userUuid, userName, userName, expiration)
	return token
}

func stringNotEqual(t *testing.T, got, expect string) {
	if got == expect {
		t.Errorf(`Got %v.`, got)
	}
}

func intEqual(t *testing.T, got, expect int) {
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func responseIsSuccess(t *testing.T, response *responseContent) {
	if response != nil {
		if response.Status != SuccessStatus {
			t.Errorf(`Response status isn't success: %v.`, response)
		}
	} else {
		t.Error(`No response.`)
	}
}

func responseIsFailure(t *testing.T, response *responseContent) {
	if response != nil {
		if response.Status != FailedStatus {
			t.Errorf(`Response status isn't failure: %v.`, response)
		}
	} else {
		t.Error(`No response.`)
	}
}

func jsonStringContains(t *testing.T, got []byte, expect string) {
	if got != nil {
		gotString := utils.CleanupStrings(string(got))
		expectString := utils.CleanupStrings(expect)
		if !strings.Contains(gotString, expectString) {
			t.Errorf(`Response doesn't include %v : %v`, expectString, gotString)
		}
	} else {
		t.Error(`Response is nil.`)
	}
}

func jsonStringDoesntContain(t *testing.T, got []byte, expect string) {
	if got != nil {
		gotJson := utils.CleanupStrings(string(got))
		expectJson := utils.CleanupStrings(expect)
		if strings.Contains(gotJson, expectJson) {
			t.Errorf(`Response includes %v : %v`, expectJson, gotJson)
		}
	} else {
		t.Error(`Response is nil.`)
	}
}

func setDefaultTestSleepTimeAndTimeOutThreshold() {
	database.ForceTestSleepTimeAndTimeOutThreshold("test", 0, 5000)
}

func RunTestConnectDbServersIncorrect(t *testing.T) {
	configuration.LoadConfiguration("abc/configuration.yml")
	dbName := "alpha"
	_, err := database.GetDbByName(dbName)
	if err == nil {
		t.Errorf(`Database %q found!`, dbName)
	}
}

func RunTestConnectDbServers(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	dbName := "test"
	_, err := database.GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func RunTestRecreateDb(t *testing.T) {
	dbName := "test"
	db, err := database.GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := database.RecreateDb(context.Background(), db, dbName); err != nil {
		t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
	}
}

func stringToJson(in string) GridPost {
	var post GridPost
	_ = json.Unmarshal([]byte(in), &post)
	return post
}
