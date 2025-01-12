// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"os"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type CustomRoundRobin struct{}

func (b *CustomRoundRobin) Balance(message kafka.Message, partitions ...int) int {
	for _, header := range message.Headers {
		if header.Key == "gridUuid" {
			gridUuid := string(header.Value)
			hash := calcHash(gridUuid)
			nbPartitions := uint32(len(partitions))
			balance := int(hash % nbPartitions)
			return balance
		}
	}
	return 0
}

func calcHash(input string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(input))
	return algorithm.Sum32()
}

func WriteMessage(dbName string, userUuid string, user string, gridUuid string,
	contextUuid string, requestInitiatedOn string, receivedOn string, key string, response responseContent) {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-responses"
	hostname, _ := os.Hostname()
	responseInitiatedOn := time.Now().UTC().Format(time.RFC3339Nano)

	w := kafka.Writer{
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
		{Key: "from", Value: []byte("εncooη backend")},
		{Key: "hostName", Value: []byte(hostname)},
		{Key: "dbName", Value: []byte(dbName)},
		{Key: "userUuid", Value: []byte(userUuid)},
		{Key: "user", Value: []byte(user)},
		{Key: "gridUuid", Value: []byte(gridUuid)},
		{Key: "contextUuid", Value: []byte(contextUuid)},
		{Key: "requestInitiatedOn", Value: []byte(requestInitiatedOn)},
		{Key: "requestReceivedOn", Value: []byte(receivedOn)},
		{Key: "responseInitiatedOn", Value: []byte(responseInitiatedOn)},
	}

	responseEncoded, err := json.Marshal(response)
	if err != nil {
		configuration.LogError(dbName, user, "Error marshal response:", err)
		return
	}
	configuration.Log(dbName, user, "PUSH message (%d bytes), topic: %s, key: %s, action: %s, status: %s", len(responseEncoded), topic, key, response.Action, response.Status)
	err = w.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   responseEncoded,
			Headers: headers,
		},
	)

	if err != nil {
		configuration.LogError(dbName, user, "failed to write messages:", err)
	}
	if err := w.Close(); err != nil {
		configuration.LogError(dbName, user, "failed to close writer:", err)
	}
}
