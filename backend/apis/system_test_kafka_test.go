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
