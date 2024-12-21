// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"strings"

	"d.lambert.fr/encoon/configuration"
	"github.com/segmentio/kafka-go"
)

func SetAndStartKafkaReader() {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.Topic
	groupID := configuration.GetConfiguration().Kafka.GroupID

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(kafkaBrokers, ","),
		Topic:   topic,
		GroupID: groupID,
	})

	configuration.Log("", "", "Connected to Kafka topic %s through brokers %s.", topic, kafkaBrokers)
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			configuration.LogError("", "", "could not read message: %v", err)
			continue
		}

		configuration.Log("", "", "Message from topic %s, partition %d and offset %d: key %s value %s", m.Topic, m.Partition, m.Offset, m.Key, m.Value)
	}
}
