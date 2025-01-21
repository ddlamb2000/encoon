// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"d.lambert.fr/encoon/configuration"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
)

var consumers = struct {
	sync.RWMutex
	m map[string]*kafka.Reader
}{m: make(map[string]*kafka.Reader)}

var consumerShutdown = struct {
	sync.RWMutex
	shutdown map[string]bool
}{shutdown: make(map[string]bool)}

func ReadMessagesFromKafka() {
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		go readMessages(dbConfig.Name)
	}
}

func ShutdownKafkaConsumers() {
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		consumer := consumers.m[dbConfig.Name]
		if consumer != nil {
			configuration.Log(dbConfig.Name, "", "Stop requested for Kafka consumer")
			consumerShutdown.shutdown[dbConfig.Name] = true
			consumer.Close()
		}
	}
}

func getConsumer(dbName string) (*kafka.Reader, error) {
	if dbName == "" || dbName == "undefined" {
		return nil, configuration.LogAndReturnError("", "", "Missing database name parameter.")
	}
	consumers.RLock()
	consumer := consumers.m[dbName]
	consumers.RUnlock()
	if consumer != nil {
		return consumer, nil
	}
	consumers.Lock()
	defer consumers.Unlock()
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-requests"
	groupID := configuration.GetConfiguration().Kafka.GroupID + "-" + dbName
	configuration.Log(dbName, "", "Creating Kafka consumer")
	consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:          strings.Split(kafkaBrokers, ","),
		Topic:            topic,
		GroupID:          groupID,
		MaxBytes:         10e3,
		MaxWait:          10 * time.Millisecond,
		RebalanceTimeout: 10 * time.Second,
		CommitInterval:   time.Second,
	})
	consumers.m[dbName] = consumer
	consumerShutdown.shutdown[dbName] = false
	return consumer, nil
}

func readMessages(dbName string) {
	consumer, err := getConsumer(dbName)
	if err != nil {
		configuration.LogError(dbName, "", "Error getting Kafka consumer", err)
		return
	}
	configuration.Log(dbName, "", "Read messages")
	for {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			if consumerShutdown.shutdown[dbName] {
				configuration.Log(dbName, "", "Kafka consumer stopped")
				return
			}
			configuration.LogError(dbName, "", "Error reading message", err)
			continue
		}
		go handleMessage(dbName, message)
	}
}

func handleMessage(dbName string, message kafka.Message) {
	requestReceivedOn := time.Now().UTC().Format(time.RFC3339Nano)
	requestInitiatedOn, tokenString, gridUuid, contextUuid := getDataFromHeaders(message)
	messageKey := string(message.Key)
	var content requestContent
	var response responseContent
	user := ""
	userUuid := ""
	if err := json.Unmarshal(message.Value, &content); err != nil {
		configuration.LogError(dbName, "", "Error message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, action: %s %s", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, content.Action, content.ActionText, err)
		response = invalidMessage(content)
	} else {
		configuration.Log(dbName, "", "PULL Message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, action: %s %s", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, content.Action, content.ActionText)
		if content.Action == ActionHeartbeat {
			response = heartBeat(content)
		} else if content.Action == ActionAuthentication {
			response = handleAuthentication(dbName, content)
		} else if content.Action == ActionLogout {
			response = logOut(content)
		} else {
			userUuid, user, response = validMessage(messageKey, dbName, tokenString, content)
		}
	}
	WriteMessage(dbName, userUuid, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, string(message.Key), response)
}

func getDataFromHeaders(message kafka.Message) (string, string, string, string) {
	requestInitiatedOn := ""
	tokenString := ""
	gridUuid := ""
	contextUuid := ""
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
	return requestInitiatedOn, tokenString, gridUuid, contextUuid
}

func validMessage(messageKey string, dbName string, tokenString string, content requestContent) (string, string, responseContent) {
	token, err := jwt.Parse(tokenString, getTokenParsingHandler(dbName))
	if err != nil {
		return "", "", noToken(messageKey, dbName, content)
	} else {
		if token == nil {
			return "", "", notAuthorization(messageKey, dbName, content)
		} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userUuid, user, tokenExpired := getDataFromJWTClaims(claims)
			if tokenExpired {
				return "", "", expired(messageKey, dbName, user, content)
			} else {
				return userUuid, user, handleActions(dbName, userUuid, user, content)
			}
		} else {
			return "", "", invalidToken(messageKey, dbName, content)
		}
	}
}

func getDataFromJWTClaims(claims jwt.MapClaims) (string, string, bool) {
	today := time.Now()
	expiration := claims["expires"]
	expirationDate, _ := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", expiration))
	userUuid := fmt.Sprintf("%v", claims["userUuid"])
	user := fmt.Sprintf("%v", claims["user"])
	return userUuid, user, today.After(expirationDate)
}

func invalidMessage(content requestContent) responseContent {
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		TextMessage: "Incorrect message",
	}
}

func handleActions(dbName string, userUuid string, userName string, content requestContent) responseContent {
	if content.Action == ActionLoad {
		return executeActionGrid(dbName, userUuid, userName, content, GetGridsRows)
	} else if content.Action == ActionChangeGrid {
		return executeActionGrid(dbName, userUuid, userName, content, PostGridsRows)
	} else if content.Action == ActionLocateGrid {
		return locate(content)
	} else {
		return invalidAction(dbName, content)
	}
}

func heartBeat(content requestContent) responseContent {
	return responseContent{
		Status: SuccessStatus,
		Action: content.Action,
	}
}

func logOut(content requestContent) responseContent {
	return responseContent{
		Status:      SuccessStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		TextMessage: "User logged out",
	}
}

func locate(content requestContent) responseContent {
	return responseContent{
		Status:     SuccessStatus,
		Action:     content.Action,
		ActionText: content.ActionText,
		GridUuid:   content.GridUuid,
		ColumnUuid: content.ColumnUuid,
		Uuid:       content.Uuid,
	}
}

func notAuthorization(messageKey, dbName string, content requestContent) responseContent {
	configuration.LogError(dbName, "", "No authorization for message %s action: %s %s", messageKey, content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		GridUuid:    content.GridUuid,
		TextMessage: "No authorization",
	}
}

func expired(messageKey, dbName string, userName string, content requestContent) responseContent {
	configuration.LogError(dbName, userName, "Authorization expired for message %s action: %s %s", messageKey, content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		GridUuid:    content.GridUuid,
		TextMessage: "Authorization expired",
	}
}

func invalidAction(dbName string, content requestContent) responseContent {
	configuration.Log(dbName, "", "Invalid action: %s %s", content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		GridUuid:    content.GridUuid,
		TextMessage: "Invalid action (" + content.Action + ")",
	}
}

func noToken(messageKey, dbName string, content requestContent) responseContent {
	configuration.LogError(dbName, "", "No token for message %s action: %s %s", messageKey, content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		GridUuid:    content.GridUuid,
		TextMessage: "Missing authorization",
	}
}

func invalidToken(messageKey, dbName string, content requestContent) responseContent {
	configuration.LogError(dbName, "", "Invalid token for message %s action: %s %s", messageKey, content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		GridUuid:    content.GridUuid,
		TextMessage: "Invalid token",
	}
}

type ActionGridDataFunc func(ct context.Context, uri string, p ApiParameters, payload GridPost) GridResponse

func executeActionGrid(dbName string, userUuid string, userName string, content requestContent, f ActionGridDataFunc) responseContent {
	parameters := getRequestParameters(dbName, userUuid, userName, content)
	response := f(context.Background(), "", parameters, content.DataSet)
	if response.Err != nil {
		return responseContent{
			Status:      FailedStatus,
			Action:      content.Action,
			ActionText:  content.ActionText,
			GridUuid:    content.GridUuid,
			Uuid:        content.Uuid,
			TextMessage: response.Err.Error(),
		}
	}
	return responseContent{
		Status:     SuccessStatus,
		Action:     content.Action,
		ActionText: content.ActionText,
		GridUuid:   content.GridUuid,
		Uuid:       content.Uuid,
		DataSet:    response,
	}
}
