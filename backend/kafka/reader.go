// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

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
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbConfig.Name + "-requests"
		groupID := configuration.GetConfiguration().Kafka.GroupID + "-" + dbConfig.Name
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
		CommitInterval:   time.Second,
	})

	configuration.Log(dbName, "", "Read messages on topic %s through brokers %s with consumer group %s.", topic, kafkaBrokers, groupID)
	for {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			configuration.LogError(dbName, "", "could not read message: %v", err)
			continue
		}
		handleMessage(dbName, message)
	}
}

func handleMessage(dbName string, message kafka.Message) {
	requestInitiatedOn := ""
	tokenString := ""
	gridUuid := ""
	contextUuid := ""
	requestReceivedOn := time.Now().UTC().Format(time.RFC3339Nano)
	for _, header := range message.Headers {
		switch header.Key {
		case "requestInitiatedOn":
			requestInitiatedOn = string(header.Value)
		case "jwt":
			tokenString = string(header.Value)
		case "gridUuid":
			gridUuid = string(header.Value)
		case "contextUuid":
			contextUuid = string(header.Value)
		}
	}
	var content requestContent
	var response responseContent
	user := ""
	userUuid := ""
	if err := json.Unmarshal(message.Value, &content); err != nil {
		configuration.LogError(dbName, "", "Error message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText, err)
		response = responseContent{
			Status:      FailedStatus,
			Action:      content.Action,
			ActionText:  content.ActionText,
			TextMessage: "Incorrect message",
		}
	} else {
		if content.Action == ActionHeartbeat {
			response = responseContent{
				Status: SuccessStatus,
				Action: content.Action,
			}
		} else if content.Action == ActionAuthentication {
			configuration.Log(dbName, "", "PULL Message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
			response = authentication(dbName, content)
		} else if content.Action == ActionLogout {
			configuration.Log(dbName, "", "PULL Message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
			response = responseContent{
				Status:      SuccessStatus,
				Action:      content.Action,
				ActionText:  content.ActionText,
				TextMessage: "User logged out",
			}
		} else {
			token, err := jwt.Parse(tokenString, getTokenParsingHandler(dbName))
			if err != nil {
				configuration.LogError(dbName, "", "Invalid token for message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText, err)
				response = responseContent{
					Status:      FailedStatus,
					Action:      content.Action,
					ActionText:  content.ActionText,
					TextMessage: "Invalid token",
				}
			} else {
				if token == nil {
					configuration.LogError(dbName, "", "No authorization for message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
					response = responseContent{
						Status:      FailedStatus,
						Action:      content.Action,
						ActionText:  content.ActionText,
						TextMessage: "No authorization",
					}
				} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					today := time.Now()
					expiration := claims["expires"]
					expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))
					user = fmt.Sprintf("%v", claims["user"])
					userUuid = fmt.Sprintf("%v", claims["userUuid"])
					if today.After(expirationDate) {
						configuration.LogError(dbName, user, "Authorization expired for message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
						response = responseContent{
							Status:      FailedStatus,
							Action:      content.Action,
							ActionText:  content.ActionText,
							TextMessage: "Authorization expired",
						}
					} else {
						configuration.Log(dbName, user, "PULL Message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
						if content.Action == ActionLoad {
							response = getGrid(dbName, userUuid, user, content)
						} else if content.Action == ActionChangeGrid {
							response = postGrid(dbName, userUuid, user, content)
						} else if content.Action == ActionLocateGrid {
							response = locate(dbName, content)
						} else {
							configuration.Log(dbName, "", "Invalid action: % %ss", content.Action, content.ActionText)
							response = responseContent{
								Status:      FailedStatus,
								Action:      content.Action,
								ActionText:  content.ActionText,
								TextMessage: "Invalid action (" + content.Action + ")",
							}
						}
					}
				} else {
					configuration.LogError(dbName, "", "Invalid request for message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
					response = responseContent{
						Status:      FailedStatus,
						Action:      content.Action,
						ActionText:  content.ActionText,
						TextMessage: "Invalid request",
					}
				}
			}
		}
	}
	WriteMessage(dbName, userUuid, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, string(message.Key), response)
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
