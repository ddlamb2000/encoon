// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"hash/fnv"
	"os"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type RoundRobin struct {
}

func (b *RoundRobin) Balance(msg kafka.Message, partitions ...int) (partition int) {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(msg.Key))
	nbPartitions := uint32(len(partitions))
	hash := algorithm.Sum32()
	balance := int(hash % nbPartitions)
	return balance
}

func WriteMessage(dbName string, userUuid string, user string, requestKey []byte, requestInitiatedOn []byte, receivedOn []byte, response []byte) {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-responses"
	hostname, _ := os.Hostname()
	responseInitiatedOn := []byte(time.Now().UTC().Format(time.RFC3339Nano))

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
		Balancer:               &RoundRobin{},
	}

	key := utils.GetNewUUID()
	headers := []kafka.Header{
		{Key: "from", Value: []byte("εncooη backend")},
		{Key: "hostName", Value: []byte(hostname)},
		{Key: "dbName", Value: []byte(dbName)},
		{Key: "userUuid", Value: []byte(userUuid)},
		{Key: "user", Value: []byte(user)},
		{Key: "requestKey", Value: requestKey},
		{Key: "requestInitiatedOn", Value: requestInitiatedOn},
		{Key: "requestReceivedOn", Value: receivedOn},
		{Key: "responseInitiatedOn", Value: responseInitiatedOn},
	}
	configuration.Log(dbName, "", "{PUSH} %d bytes, topic: %s, key: %s", len(response), topic, key)
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
