// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

const iterations = 10
const interval = time.Millisecond * 200

func waitAndRepeat(t *testing.T, consumer *kafka.Reader, key string, f func(t *testing.T, consumer *kafka.Reader, key string) (*responseContent, bool)) *responseContent {
	for i := 0; i < iterations; i++ {
		response, stop := f(t, consumer, key)
		if stop {
			return response
		}
		time.Sleep(interval)
	}
	return nil
}

func readKafkaTestMessage(t *testing.T, consumer *kafka.Reader, key string) (*responseContent, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	message, err := consumer.ReadMessage(ctx)
	if err != nil {
		t.Errorf("Error reading message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, %v", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, err)
		return nil, false
	}
	t.Logf("Got message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset)
	if err := consumer.CommitMessages(ctx, message); err != nil {
		t.Errorf("Error committing message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, %v", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, err)
		return nil, false
	}
	if string(message.Key) == key {
		var response responseContent
		if err := json.Unmarshal(message.Value, &response); err != nil {
			t.Errorf("Error unmarshal message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, %v", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, err)
			return nil, false
		}
		return &response, true
	}
	return nil, false
}

func getKafkaTestConsumer(dbName string) *kafka.Reader {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-responses"
	groupID := configuration.GetConfiguration().Kafka.GroupID + "-" + dbName
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          strings.Split(kafkaBrokers, ","),
		Topic:            topic,
		GroupID:          groupID,
		MaxBytes:         10e3,
		MaxWait:          10 * time.Millisecond,
		RebalanceTimeout: 2 * time.Second,
		CommitInterval:   time.Second,
	})
	return consumer
}

func RunSystemTestKafka(t *testing.T) {
	configuration.LoadConfiguration("../testData/systemTest.yml")
	go ReadMessagesFromKafka()
	consumer := getKafkaTestConsumer("test")

	t.Run("Base test using Kafka", func(t *testing.T) {
		expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
		token, _ := getNewToken("test", "root", model.UuidRootUser, "root", "root", expiration)
		key := "key5"
		message := responseContent{
			Action: ActionHeartbeat,
		}
		writeTestMessage(t, "test", "rootuuid", "root", token, model.UuidUsers, "context1", key, message)
		response := waitAndRepeat(t, consumer, key, readKafkaTestMessage)
		if response.Status != SuccessStatus {
			t.Errorf(`Response status isn't success: %v.`, response)
		}
	})
	consumer.Close()
	ShutdownKafkaProducers()
	ShutdownKafkaConsumers()
}

func writeTestMessage(t *testing.T, dbName string, userUuid string, user string, token string, gridUuid string,
	contextUuid string, key string, message responseContent) {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-requests"
	requestInitiatedOn := time.Now().UTC().Format(time.RFC3339Nano)

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(kafkaBrokers, ",")[:]...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		MaxAttempts:            3,
		WriteBackoffMax:        5 * time.Millisecond,
		BatchSize:              10,
		BatchTimeout:           50 * time.Millisecond,
		RequiredAcks:           -1,
		Compression:            compress.Gzip,
		Balancer:               &CustomRoundRobin{},
	}

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
	t.Logf("PUSH message (%d bytes), topic: %s, key: %s, action: %s, status: %s", len(messageEncoded), topic, key, message.Action, message.Status)
	err := (*writer).WriteMessages(
		context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   messageEncoded,
			Headers: headers,
		},
	)

	if err != nil {
		t.Error("failed to write messages:", err)
	}

	if err := writer.Close(); err != nil {
		t.Error("Failed to close Kafka producer:", err)
	}
}
