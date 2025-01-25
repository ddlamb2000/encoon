// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

var kafakMessageNumber = 0
var kafkaTestConsumer *kafka.Reader = nil
var kafkaTestProducer *kafka.Writer = nil
var kafkaBaddbConsumer *kafka.Reader = nil
var kafkaBaddbProducer *kafka.Writer = nil
var contextUuid = utils.GetNewUUID()

func RunSystemTestKafka(t *testing.T) {
	t.Run("GetConsumer", func(t *testing.T) {
		consumer, err := getConsumer("test")
		if err != nil || consumer == nil {
			t.Errorf(`No consumer: %v.`, err)
		}
	})

	t.Run("GetConsumerNoDb", func(t *testing.T) {
		consumer, err := getConsumer("")
		expect := "Missing database name parameter."
		if err == nil || err.Error() != expect {
			t.Errorf("Got err %v instead of %v.", err, expect)
			consumer.Close()
		}
	})

	t.Run("GetProducer", func(t *testing.T) {
		producer, err := getProducer("test")
		if err != nil || producer == nil {
			t.Errorf(`No consumer: %v.`, err)
		}
	})

	t.Run("GetProducerNoDb", func(t *testing.T) {
		producer, err := getProducer("")
		expect := "Missing database name parameter."
		if err == nil || err.Error() != expect {
			t.Errorf("Got err %v instead of %v.", err, expect)
			producer.Close()
		}
	})

	t.Run("ReadMessagesNoDbName", func(t *testing.T) {
		err := readMessages("")
		expect := "Error getting Kafka consumer: Missing database name parameter."
		if err == nil || err.Error() != expect {
			t.Errorf("Got err %v instead of %v.", err, expect)
		}
	})

	t.Run("WriteMessagesNoDbName", func(t *testing.T) {
		err := WriteMessage("", "", "", "", "", "", "", "", responseContent{})
		expect := "Error getting Kafka producer: Missing database name parameter."
		if err == nil || err.Error() != expect {
			t.Errorf("Got err %v instead of %v.", err, expect)
		}
	})

	t.Run("WriteMessagesDefectJson", func(t *testing.T) {
		jsonMarshalImpl := jsonMarshal
		jsonMarshal = func(v any) ([]byte, error) {
			return nil, errors.New("xxx")
		} // mock function
		err := WriteMessage("test", "", "", "", "", "", "", "", responseContent{})
		expect := "Error marshal response: xxx"
		if err == nil || err.Error() != expect {
			t.Errorf("Got err %v instead of %v.", err, expect)
		}
		jsonMarshal = jsonMarshalImpl
	})

	t.Run("WriteMessagesDefectWrite", func(t *testing.T) {
		writeMessageImpl := writeMessage
		writeMessage = func(w *kafka.Writer, ctx context.Context, msgs ...kafka.Message) error {
			return errors.New("xxx")
		} // mock function
		err := WriteMessage("test", "", "", "", "", "", "", "KEY", responseContent{
			Action: ActionHeartbeat,
			Status: FailedStatus,
		})
		expect := "Failed to PUSH message (53 bytes), key: KEY, action: HEARTBEAT, status: FAILED"
		if err == nil || err.Error() != expect {
			t.Errorf("Got err %v instead of %v.", err, expect)
		}
		writeMessage = writeMessageImpl
	})

	t.Run("Heartbeat", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "", ApiParameters{
			Action: ActionHeartbeat,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"action":"HEARTBEAT","status":"SUCCESS"`)
	})

	t.Run("HeartbeatNoReadMessage", func(t *testing.T) {
		readMessageImpl := readMessage
		readMessage = func(r *kafka.Reader, ctx context.Context) (kafka.Message, error) {
			return kafka.Message{}, errors.New("xxx")
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "", ApiParameters{
			Action: ActionHeartbeat,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"action":"HEARTBEAT","status":"SUCCESS"`)
		readMessage = readMessageImpl
	})

	t.Run("HeartbeatBadjsonUnmarshal", func(t *testing.T) {
		jsonUnmarshalImpl := jsonUnmarshal
		jsonUnmarshal = func(data []byte, v any) error {
			return errors.New("xxx")
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "", ApiParameters{
			Action: ActionHeartbeat,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Incorrect message"`)
		jsonUnmarshal = jsonUnmarshalImpl
	})

	t.Run("HeartbeatBadjsonUnmarshal", func(t *testing.T) {
		jwtParseImpl := jwtParse
		jwtParse = func(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
			return nil, errors.New("yyy")
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Invalid authorization"`)
		jwtParse = jwtParseImpl
	})

	t.Run("HeartbeatNoAuthorization", func(t *testing.T) {
		jwtParseImpl := jwtParse
		jwtParse = func(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
			return nil, nil
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"No authorization"`)
		jwtParse = jwtParseImpl
	})

	t.Run("HeartbeatBadJwtClaims", func(t *testing.T) {
		getJwtClaimsImpl := getJwtClaims
		getJwtClaims = func(token *jwt.Token) jwt.Claims {
			return nil
		} // mock function
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Invalid token"`)
		getJwtClaims = getJwtClaimsImpl
	})

	t.Run("Locate", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   ActionLocateGrid,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"action":"LOCATE","status":"SUCCESS"`)
	})

	t.Run("InvalidAction", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, ApiParameters{
			Action:   "X",
			GridUuid: model.UuidGrids,
		})
		responseIsFailure(t, response)
		jsonStringContains(t, responseData, `"textMessage":"Invalid action (X)`)
	})
}

func startReadingTestMessages() {
	StartReadingMessages()
	kafkaTestProducer, kafkaTestConsumer = createKafkaTestProducer("test"), createKafkaTestConsumer("test")
	kafkaBaddbProducer, kafkaBaddbConsumer = createKafkaTestProducer("baddb"), createKafkaTestConsumer("baddb")
}

func stopReadingTestMessages() {
	kafkaTestProducer.Close()
	kafkaTestConsumer.Close()
	kafkaBaddbProducer.Close()
	kafkaBaddbConsumer.Close()
	StopReadingMessages()
}

func readKafkaTestMessage(t *testing.T, consumer *kafka.Reader, key string) (*responseContent, []byte) {
	for i := 0; i < 50; i++ {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			t.Errorf("Error reading message: %v", err)
			return nil, nil
		}
		if string(message.Key) == key {
			var response responseContent
			if err := json.Unmarshal(message.Value, &response); err != nil {
				t.Errorf("Error unmarshal message key: %s, %v", message.Key, err)
				return nil, nil
			}
			return &response, message.Value
		}
	}
	t.Errorf("No message")
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
		WriteBackoffMax:        10 * time.Millisecond,
		BatchSize:              10,
		BatchTimeout:           50 * time.Millisecond,
		RequiredAcks:           -1,
		Compression:            compress.Gzip,
		Balancer:               &CustomRoundRobin{},
	}
}

func runKafkaTestAuthRequest(t *testing.T, dbName string, message ApiParameters) (*responseContent, []byte) {
	kafakMessageNumber++
	key := fmt.Sprintf("%d-%s", kafakMessageNumber, contextUuid)
	writeTestAuthMessage(t, dbName, key, message)
	consumer := kafkaTestConsumer
	if dbName == "baddb" {
		consumer = kafkaBaddbConsumer
	}
	return readKafkaTestMessage(t, consumer, key)
}

func writeTestAuthMessage(t *testing.T, dbName, key string, message ApiParameters) {
	messageEncoded, _ := json.Marshal(message)
	producer := kafkaTestProducer
	if dbName == "baddb" {
		producer = kafkaBaddbProducer
	}
	err := (*producer).WriteMessages(
		context.Background(),
		kafka.Message{
			Key: []byte(key),
			Headers: []kafka.Header{
				{Key: "from", Value: []byte("εncooη testing")},
				{Key: "url", Value: []byte("http://localhost")},
				{Key: "contextUuid", Value: []byte(contextUuid)},
				{Key: "dbName", Value: []byte(dbName)},
				{Key: "requestInitiatedOn", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
			},
			Value: messageEncoded,
		},
	)

	if err != nil {
		t.Error("Failed to write messages:", err)
	}
}

func getTokenForUser(dbName, userName, userUuid string) string {
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().JwtExpiration) * time.Minute)
	token, _ := getNewToken(dbName, userName, userUuid, userName, userName, expiration)
	return token
}

func runKafkaTestRequest(t *testing.T, dbName, userName, userUuid, gridUuid string, message ApiParameters) (*responseContent, []byte) {
	token := getTokenForUser(dbName, userName, userUuid)
	kafakMessageNumber++
	key := fmt.Sprintf("%d-%s", kafakMessageNumber, contextUuid)
	writeTestMessage(t, dbName, userUuid, userName, token, gridUuid, key, message)
	consumer := kafkaTestConsumer
	if dbName == "baddb" {
		consumer = kafkaBaddbConsumer
	}
	return readKafkaTestMessage(t, consumer, key)
}

func runKafkaTestRequestWithToken(t *testing.T, dbName, userName, userUuid, gridUuid, token string, message ApiParameters) (*responseContent, []byte) {
	kafakMessageNumber++
	key := fmt.Sprintf("%d-%s", kafakMessageNumber, contextUuid)
	writeTestMessage(t, dbName, userUuid, userName, token, gridUuid, key, message)
	consumer := kafkaTestConsumer
	if dbName == "baddb" {
		consumer = kafkaBaddbConsumer
	}
	return readKafkaTestMessage(t, consumer, key)
}

func writeTestMessage(t *testing.T, dbName, userUuid, user, token, gridUuid, key string, message ApiParameters) {
	messageEncoded, _ := json.Marshal(message)
	producer := kafkaTestProducer
	if dbName == "baddb" {
		producer = kafkaBaddbProducer
	}
	err := (*producer).WriteMessages(
		context.Background(),
		kafka.Message{
			Key: []byte(key),
			Headers: []kafka.Header{
				{Key: "from", Value: []byte("εncooη testing")},
				{Key: "url", Value: []byte("http://localhost")},
				{Key: "contextUuid", Value: []byte(contextUuid)},
				{Key: "dbName", Value: []byte(dbName)},
				{Key: "userUuid", Value: []byte(userUuid)},
				{Key: "user", Value: []byte(user)},
				{Key: "jwt", Value: []byte(token)},
				{Key: "gridUuid", Value: []byte(gridUuid)},
				{Key: "requestInitiatedOn", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
			},
			Value: messageEncoded,
		},
	)

	if err != nil {
		t.Error("Failed to write messages:", err)
	}
}
