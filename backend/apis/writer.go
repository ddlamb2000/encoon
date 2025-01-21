// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"os"
	"strings"
	"sync"
	"time"

	"d.lambert.fr/encoon/configuration"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

var producers = struct {
	sync.RWMutex
	m map[string]*kafka.Writer
}{m: make(map[string]*kafka.Writer)}

var producerShutdown = struct {
	sync.RWMutex
	shutdown map[string]bool
}{shutdown: make(map[string]bool)}

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

func ShutdownKafkaProducers() {
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		writer := producers.m[dbConfig.Name]
		if writer != nil {
			configuration.Log(dbConfig.Name, "", "Stop requested for Kafka producer")
			producerShutdown.shutdown[dbConfig.Name] = true
			if err := writer.Close(); err != nil {
				configuration.LogError(dbConfig.Name, "", "Failed to close Kafka producer:", err)
			} else {
				configuration.Log(dbConfig.Name, "", "Kafka producer stopped")
			}
		}
	}
}

func calcHash(input string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(input))
	return algorithm.Sum32()
}

func getProducer(dbName string) (*kafka.Writer, error) {
	if dbName == "" || dbName == "undefined" {
		return nil, configuration.LogAndReturnError("", "", "Missing database name parameter.")
	}
	producers.RLock()
	writer := producers.m[dbName]
	producers.RUnlock()
	if writer != nil {
		return writer, nil
	}
	producers.Lock()
	defer producers.Unlock()
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-responses"
	configuration.Log(dbName, "", "Creating Kafka producer")
	writer = &kafka.Writer{
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
	producers.m[dbName] = writer
	producerShutdown.shutdown[dbName] = false
	return writer, nil
}

func WriteMessage(dbName string, userUuid string, user string, gridUuid string,
	contextUuid string, requestInitiatedOn string, receivedOn string, key string, message responseContent) {
	hostname, _ := os.Hostname()
	responseInitiatedOn := time.Now().UTC().Format(time.RFC3339Nano)

	writer, err := getProducer(dbName)
	if err != nil {
		configuration.LogError(dbName, "", "Error getting Kafka producer", err)
		return
	}
	if producerShutdown.shutdown[dbName] {
		configuration.Log(dbName, user, "Message not sent (producer is shuting down), key: %s, action: %s, status: %s", key, message.Action, message.Status)
		return
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

	messageEncoded, err := json.Marshal(message)
	if err != nil {
		configuration.LogError(dbName, user, "Error marshal response:", err)
		return
	}
	err = (*writer).WriteMessages(
		context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   messageEncoded,
			Headers: headers,
		},
	)
	if err != nil {
		configuration.LogError(dbName, user, "Failed to PUSH message (%d bytes), key: %s, action: %s, status: %s", len(messageEncoded), key, message.Action, message.Status)
	} else {
		configuration.Log(dbName, user, "PUSH message (%d bytes), key: %s, action: %s, status: %s", len(messageEncoded), key, message.Action, message.Status)
	}
}
