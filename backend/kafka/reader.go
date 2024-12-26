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

type messageUserText struct {
	UserText string `json:"userText"`
}

type requestContent struct {
	Action   string `json:"action"`
	GridUuid string `json:"griduuid,omitempty"`
	RowUuid  string `json:"rowuuid,omitempty"`
}

type responseContent struct {
	Status   string `json:"status"`
	Action   string `json:"action"`
	GridUuid string `json:"griduuid,omitempty"`
	RowUuid  string `json:"rowuuid,omitempty"`
}

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

			var text messageUserText
			var content requestContent
			if err = json.Unmarshal(m.Value, &text); err != nil {
				configuration.LogError("", "", "error unmarshal m.Value:", err)
				continue
			}
			if err = json.Unmarshal([]byte(text.UserText), &content); err != nil {
				configuration.LogError("", "", "error unmarshal messageText.UserText:", err)
				continue
			}

			configuration.Log("", "", "content: %s", content)

			response := responseContent{
				Status:   "OK",
				Action:   content.Action,
				GridUuid: content.GridUuid,
				RowUuid:  content.RowUuid,
			}

			responseEncoded, err := json.Marshal(response)
			if err != nil {
				configuration.LogError("", "", "error marshal response:", err)
				continue
			}

			responseMessage := messageUserText{
				UserText: string(responseEncoded),
			}
			messageUserTextEncoded, err := json.Marshal(responseMessage)
			if err != nil {
				configuration.LogError("", "", "error marshal responseMessage:", err)
				continue
			}

			WriteMessage(m.Key, messageUserTextEncoded)
		}
	}

	if err := r.Close(); err != nil {
		configuration.LogError("", "", "failed to close reader:", err)
	}
}
