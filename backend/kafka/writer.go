// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"os"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type CustomRoundRobin struct{}

func (b *CustomRoundRobin) Balance(message kafka.Message, partitions ...int) (partition int) {
	gridUuid := ""
	for _, header := range message.Headers {
		if header.Key == "gridUuid" {
			gridUuid = string(header.Value)
			break
		}
	}
	algorithm := fnv.New32a()
	algorithm.Write([]byte(gridUuid))
	nbPartitions := uint32(len(partitions))
	hash := algorithm.Sum32()
	balance := int(hash % nbPartitions)
	return balance
}

func WriteMessage(dbName string, userUuid string, user string, gridUuid string,
	requestInitiatedOn string, receivedOn string, requestKey string, response responseContent) {
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
		BatchTimeout:           100 * time.Millisecond,
		RequiredAcks:           -1,
		Compression:            compress.Gzip,
		Balancer:               &CustomRoundRobin{},
	}

	key := utils.GetNewUUID()
	headers := []kafka.Header{
		{Key: "from", Value: []byte("εncooη backend")},
		{Key: "hostName", Value: []byte(hostname)},
		{Key: "dbName", Value: []byte(dbName)},
		{Key: "userUuid", Value: []byte(userUuid)},
		{Key: "user", Value: []byte(user)},
		{Key: "requestKey", Value: []byte(requestKey)},
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
