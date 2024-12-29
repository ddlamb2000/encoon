// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

func WriteMessage(dbName string, requestKey []byte, initiatedOn []byte, response []byte) {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-responses"

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
		Balancer:               &kafka.RoundRobin{},
	}

	key := utils.GetNewUUID()
	headers := []kafka.Header{
		{Key: "from", Value: []byte("backend")},
		{Key: "requestKey", Value: requestKey},
		{Key: "initiatedOn", Value: initiatedOn},
	}
	configuration.Log(dbName, "", "{PUSH} %d bytes, topic: %s, key: %s, value: %s", len(response), topic, key, response)
	err := w.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   response,
			Headers: headers,
		},
	)

	if err != nil {
		configuration.LogError(dbName, "", "failed to write messages:", err)
	}
	if err := w.Close(); err != nil {
		configuration.LogError(dbName, "", "failed to close writer:", err)
	}
}
