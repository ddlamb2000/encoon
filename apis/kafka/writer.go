// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"strings"

	"d.lambert.fr/encoon/configuration"
	"github.com/segmentio/kafka-go"
)

func WriteMessage(payload []byte) {

	configuration.Log("", "", "WriteMessage %s.", payload)

	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.Topic

	w := kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(kafkaBrokers, ",")[:]...),
		Topic:                  topic + "-4",
		AllowAutoTopicCreation: true,
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   payload,
			Value: []byte("Hello World!"),
		},
	)

	if err != nil {
		configuration.LogError("", "", "failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		configuration.LogError("", "", "failed to close writer:", err)
	}
}
