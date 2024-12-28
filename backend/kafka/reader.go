// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"d.lambert.fr/encoon/configuration"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
)

func SetAndStartKafkaReader() {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	groupID := configuration.GetConfiguration().Kafka.GroupID
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbConfig.Name + "-requests"
		go setAndStartKafkaReaderForDatabase(dbConfig.Name, kafkaBrokers, groupID, topic)
	}
}

func setAndStartKafkaReaderForDatabase(dbName string, kafkaBrokers string, groupID string, topic string) {
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
}

func handleMessage(dbName string, message kafka.Message) {
	configuration.Log(dbName, "", "Got: topic: %s, key: %s, value: %s, headers: %s", message.Topic, message.Key, message.Value, message.Headers)
	initiatedOn := []byte("")
	tokenString := []byte("")
	for _, header := range message.Headers {
		switch header.Key {
		case "initiatedOn":
			initiatedOn = header.Value
		case "jwt":
			tokenString = header.Value
		}
	}
	var content requestContent
	if err := json.Unmarshal(message.Value, &content); err != nil {
		configuration.LogError(dbName, "", "Error unmarshal message value", err)
		return
	}
	var response responseContent
	if content.Action == ActionAuthentication {
		response = authentication(dbName, content)
	} else {
		token, err := jwt.Parse(string(tokenString), getTokenParsingHandler(dbName))
		if token == nil {
			configuration.LogError(dbName, "", "No authorization")
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user"]
			today := time.Now()
			expiration := claims["expires"]
			expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))
			userName := fmt.Sprintf("%v", user)
			if today.After(expirationDate) {
				configuration.Log(dbName, userName, "Authorization expired (%v).", expirationDate)
				return
			}
		} else {
			configuration.LogError(dbName, "", "Invalid request: %v.", err)
			return
		}
		if content.Action == ActionGetGrid {
			response = getGrid(dbName, content)
		} else {
			response = responseContent{
				Status:   FailedStatus,
				Action:   content.Action,
				GridUuid: content.GridUuid,
				Uuid:     content.Uuid,
			}
		}
	}

	responseEncoded, err := json.Marshal(response)
	if err != nil {
		configuration.LogError(dbName, "", "error marshal response:", err)
		return
	}
	WriteMessage(dbName, message.Key, initiatedOn, responseEncoded)
}

func getTokenParsingHandler(dbName string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if ok := verifyToken(token); !ok {
			return nil, configuration.LogAndReturnError(dbName, "", "Unexpect signing method: %v.", token.Header["alg"])
		}
		return []byte(configuration.GetJWTSecret(dbName)), nil
	}
}

// function is available for mocking
var verifyToken = func(token *jwt.Token) bool {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	return ok
}
