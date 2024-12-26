// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"github.com/segmentio/kafka-go"
)

func SetAndStartKafkaReader() {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-master-requests"
	groupID := configuration.GetConfiguration().Kafka.GroupID

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(kafkaBrokers, ","),
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 10e3,
		MaxWait:  10 * time.Millisecond,
	})

	configuration.Log("", "", "Read messages on topic %s through brokers %s with consumer group %s.", topic, kafkaBrokers, groupID)
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			configuration.LogError("", "", "could not read message: %v", err)
			continue
		}

		err = r.CommitMessages(context.Background(), m)
		if err != nil {
			configuration.LogError("", "", "failed to commit message from topic %s, partition %d and offset %d", m.Topic, m.Partition, m.Offset)
		} else {
			configuration.Log("", "", "Got: topic: %s, key: %s, value: %s, headers: %s", m.Topic, m.Key, m.Value, m.Headers)
			WriteMessage(m.Key, m.Value)
		}
	}

	if err := r.Close(); err != nil {
		configuration.LogError("", "", "failed to close reader:", err)
	}
}
