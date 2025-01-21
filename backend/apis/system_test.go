// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

//
// Note
// System test can be run using command line `go test apis/* -v -run TestSystem`
//

package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

var testRouter = gin.Default()
var kafakMessageNumber = 0
var kafkaTestConsumer *kafka.Reader = nil
var kafkaTestProducer *kafka.Writer = nil
var kafkaBaddbConsumer *kafka.Reader = nil
var kafkaBaddbProducer *kafka.Writer = nil

var contextUuid = utils.GetNewUUID()

func TestSystem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetApiRoutes(testRouter)
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
	// t.Run("Cache", func(t *testing.T) { RunSystemTestCache(t) })
	// t.Run("NotOwnedColumn", func(t *testing.T) { RunSystemTestNotOwnedColumn(t) })
}

func getTokenForUser(dbName, userName, userUuid string) string {
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
	token, _ := getNewToken(dbName, userName, userUuid, userName, userName, expiration)
	return token
}

func runPOSTRequestForUser(dbName, userName, userUuid, uri, body string) ([]byte, int, error) {
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, w.Code, err
}

func runGETRequestForUser(dbName, userName, userUuid, uri string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("Authorization", "Bearer "+getTokenForUser(dbName, userName, userUuid))
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	responseData, err := io.ReadAll(w.Body)
	return responseData, w.Code, err
}

func httpCodeEqual(t *testing.T, code int, expectCode int) {
	if code != expectCode {
		t.Errorf(`Response code %v instead of %v.`, code, expectCode)
	}
}

func stringNotEqual(t *testing.T, got, expect string) {
	if got == expect {
		t.Errorf(`Got %v.`, got)
	}
}

func byteEqualString(t *testing.T, got []byte, expect string) {
	gotString := string(got)
	if gotString != expect {
		t.Errorf(`Got %v instead of %v.`, gotString, expect)
	}
}

func intEqual(t *testing.T, got, expect int) {
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func errorIsNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`Error: %v.`, err)
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

func readKafkaTestMessage(t *testing.T, consumer *kafka.Reader, key string) (*responseContent, []byte) {
	for i := 0; i < 50; i++ {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			t.Errorf("Error reading message: %v", err)
			return nil, nil
		}
		t.Logf("GOT response message %s [offset %d in partition %d] (expecting %s)", message.Key, message.Offset, message.Partition, key)
		if string(message.Key) == key {
			t.Log("Message matches")
			var response responseContent
			if err := json.Unmarshal(message.Value, &response); err != nil {
				t.Errorf("Error unmarshal message key: %s, %v", message.Key, err)
				return nil, nil
			}
			return &response, message.Value
		}
	}
	return nil, nil
}

func createKafkaTestConsumer(dbName string) *kafka.Reader {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-responses"
	groupID := configuration.GetConfiguration().Kafka.GroupID + "-" + dbName
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:          strings.Split(kafkaBrokers, ","),
		Topic:            topic,
		GroupID:          groupID,
		MaxBytes:         10e3,
		MaxWait:          5 * time.Millisecond,
		RebalanceTimeout: 10 * time.Second,
		CommitInterval:   time.Second,
	})
}

func createKafkaTestProducer(dbName string) *kafka.Writer {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-requests"
	return &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(kafkaBrokers, ",")[:]...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		MaxAttempts:            3,
		WriteBackoffMax:        500 * time.Millisecond,
		BatchSize:              10,
		BatchTimeout:           50 * time.Millisecond,
		RequiredAcks:           -1,
		Compression:            compress.Gzip,
		Balancer:               &CustomRoundRobin{},
	}
}

func writeTestMessage(t *testing.T, dbName, userUuid, user, token, contextUuid, gridUuid, key string, message requestContent) {
	requestInitiatedOn := time.Now().UTC().Format(time.RFC3339Nano)
	headers := []kafka.Header{
		{Key: "from", Value: []byte("εncooη testing")},
		{Key: "url", Value: []byte("http://localhost")},
		{Key: "contextUuid", Value: []byte(contextUuid)},
		{Key: "dbName", Value: []byte(dbName)},
		{Key: "userUuid", Value: []byte(userUuid)},
		{Key: "user", Value: []byte(user)},
		{Key: "jwt", Value: []byte(token)},
		{Key: "gridUuid", Value: []byte(gridUuid)},
		{Key: "requestInitiatedOn", Value: []byte(requestInitiatedOn)},
	}
	messageEncoded, _ := json.Marshal(message)
	producer := kafkaTestProducer
	if dbName == "baddb" {
		producer = kafkaBaddbProducer
	}
	t.Logf("PUSH request message key: %s", key)
	err := (*producer).WriteMessages(
		context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   messageEncoded,
			Headers: headers,
		},
	)

	if err != nil {
		t.Error("Failed to write messages:", err)
	}
}

func runKafkaTestRequest(t *testing.T, dbName, userName, userUuid, gridUuid string, message requestContent) (*responseContent, []byte) {
	token := getTokenForUser(dbName, userName, userUuid)
	kafakMessageNumber++
	key := fmt.Sprintf("%d-%s", kafakMessageNumber, contextUuid)
	writeTestMessage(t, dbName, userUuid, userName, token, contextUuid, gridUuid, key, message)
	consumer := kafkaTestConsumer
	if dbName == "baddb" {
		consumer = kafkaBaddbConsumer
	}
	return readKafkaTestMessage(t, consumer, key)
}
