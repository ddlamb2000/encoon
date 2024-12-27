// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/utils"
	"github.com/segmentio/kafka-go"
)

type requestContent struct {
	Action   string `json:"action"`
	GridUuid string `json:"griduuid,omitempty"`
	RowUuid  string `json:"rowuuid,omitempty"`
	Dbname   string `json:"dbname,omitempty"`
	Userid   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
}

type responseContent struct {
	Status       string `json:"status"`
	Action       string `json:"action"`
	GridUuid     string `json:"griduuid,omitempty"`
	RowUuid      string `json:"rowuuid,omitempty"`
	ErrorMessage string `json:"errormessage,omitempty"`
}

func SetAndStartKafkaReader() {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-master-requests"
	groupID := configuration.GetConfiguration().Kafka.GroupID + "-" + utils.GetNewUUID()

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(kafkaBrokers, ","),
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 10e3,
		MaxWait:  10 * time.Millisecond,
	})

	configuration.Log("", "", "Read messages on topic %s through brokers %s with consumer group %s.", topic, kafkaBrokers, groupID)
	for {
		m, err := consumer.FetchMessage(context.Background())
		if err != nil {
			configuration.LogError("", "", "could not read message: %v", err)
			continue
		}

		err = consumer.CommitMessages(context.Background(), m)
		if err != nil {
			configuration.LogError("", "", "failed to commit message from topic %s, partition %d and offset %d", m.Topic, m.Partition, m.Offset)
		} else {
			configuration.Log("", "", "Got: topic: %s, key: %s, value: %s, headers: %s", m.Topic, m.Key, m.Value, m.Headers)

			initiatedOn := []byte("")
			for _, header := range m.Headers {
				if header.Key == "initiatedOn" {
					initiatedOn = header.Value
				}
			}

			var content requestContent
			if err = json.Unmarshal(m.Value, &content); err != nil {
				configuration.LogError("", "", "Error unmarshal message value", err)
				continue
			}

			var response responseContent
			if content.Action == "login" {
				configuration.Log("", "", "try to login %s, %s, %s", content.Dbname, content.Userid, content.Password)
				userUuid, firstName, lastName, timeOut, err := database.IsDbAuthorized(context.Background(), content.Dbname, content.Userid, content.Password)
				configuration.Log("", "", " %s, %s, %s", userUuid, firstName, lastName)
				if err != nil || userUuid == "" {
					if timeOut {
						response = responseContent{
							Status:       "KO",
							ErrorMessage: "Timeout",
						}

					} else {
						response = responseContent{
							Status:       "KO",
							ErrorMessage: "Invalid username or passphrase",
						}
					}
				} else {
					response = responseContent{
						Status: "OK",
					}
				}
			} else {
				response = responseContent{
					Status:   "OK",
					Action:   content.Action,
					GridUuid: content.GridUuid,
					RowUuid:  content.RowUuid,
				}
			}

			responseEncoded, err := json.Marshal(response)
			if err != nil {
				configuration.LogError("", "", "error marshal response:", err)
				continue
			}
			WriteMessage(m.Key, initiatedOn, responseEncoded)
		}
	}

	if err := consumer.Close(); err != nil {
		configuration.LogError("", "", "Failed to close reader:", err)
	}
}
