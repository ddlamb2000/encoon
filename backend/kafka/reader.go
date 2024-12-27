// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"github.com/segmentio/kafka-go"
)

func SetAndStartKafkaReader() {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	groupID := configuration.GetConfiguration().Kafka.GroupID
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbConfig.Name + "-requests"
		go SetAndStartKafkaReaderForDatabase(dbConfig.Name, kafkaBrokers, groupID, topic)
	}
}

func SetAndStartKafkaReaderForDatabase(dbName string, kafkaBrokers string, groupID string, topic string) {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          strings.Split(kafkaBrokers, ","),
		Topic:            topic,
		GroupID:          groupID,
		MaxBytes:         10e3,
		MaxWait:          10 * time.Millisecond,
		RebalanceTimeout: 2 * time.Second,
	})

	configuration.Log(dbName, "", "Read messages on topic %s through brokers %s with consumer group %s.", topic, kafkaBrokers, groupID)
	for {
		message, err := consumer.FetchMessage(context.Background())
		if err != nil {
			configuration.LogError(dbName, "", "could not read message: %v", err)
			continue
		}
		err = consumer.CommitMessages(context.Background(), message)
		if err != nil {
			configuration.LogError(dbName, "", "failed to commit message from topic %s, partition %d and offset %d", message.Topic, message.Partition, message.Offset)
		} else {
			handleMessage(dbName, message)
		}
	}

	if err := consumer.Close(); err != nil {
		configuration.LogError(dbName, "", "Failed to close reader:", err)
	}
}

func handleMessage(dbName string, message kafka.Message) {
	configuration.Log(dbName, "", "Got: topic: %s, key: %s, value: %s, headers: %s", message.Topic, message.Key, message.Value, message.Headers)

	initiatedOn := []byte("")
	for _, header := range message.Headers {
		if header.Key == "initiatedOn" {
			initiatedOn = header.Value
		}
	}

	var content requestContent
	if err := json.Unmarshal(message.Value, &content); err != nil {
		configuration.LogError(dbName, "", "Error unmarshal message value", err)
		return
	}

	var response responseContent
	if content.Action == ActionAuthentication {
		response = authentication(dbName, content.Action, content)
	} else {
		response = responseContent{
			Status:   SuccessStatus,
			Action:   content.Action,
			GridUuid: content.GridUuid,
			RowUuid:  content.RowUuid,
		}
	}

	responseEncoded, err := json.Marshal(response)
	if err != nil {
		configuration.LogError(dbName, "", "error marshal response:", err)
		return
	}
	WriteMessage(dbName, message.Key, initiatedOn, responseEncoded)
}
