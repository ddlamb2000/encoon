// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"strings"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/utils"
	"github.com/segmentio/kafka-go"
)

func WriteMessage(payload []byte) {

	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-master-responses"

	w := kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(kafkaBrokers, ",")[:]...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	key := utils.GetNewUUID()
	headers := []kafka.Header{{Key: "from", Value: []byte("backend")}}
	configuration.Log("", "", "Send: topic: %s, key: %s, value: %s, headers: %s", topic, key, payload, headers)
	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   payload,
			Headers: headers,
		},
	)

	if err != nil {
		configuration.LogError("", "", "failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		configuration.LogError("", "", "failed to close writer:", err)
	}
}
