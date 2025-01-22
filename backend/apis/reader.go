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
	"github.com/golang-jwt/jwt/v5"
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
	var request ApiParameters
	var response responseContent
	user := ""
	userUuid := ""
	if err := json.Unmarshal(message.Value, &request); err != nil {
		configuration.LogError(dbName, "", "Error message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, action: %s %s", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, request.Action, request.ActionText, err)
		response = invalidMessage(request)
	} else {
		configuration.Log(dbName, "", "PULL Message (%d bytes), topic: %s, key: %s, partition: %d, offset: %d, action: %s %s", len(message.Value), message.Topic, message.Key, message.Partition, message.Offset, request.Action, request.ActionText)
		if request.Action == ActionHeartbeat {
			response = heartBeat(request)
		} else if request.Action == ActionAuthentication {
			response = handleAuthentication(dbName, request)
		} else if request.Action == ActionLogout {
			response = logOut(request)
		} else {
			userUuid, user, response = validMessage(messageKey, dbName, tokenString, request)
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

func validMessage(messageKey string, dbName string, tokenString string, request ApiParameters) (string, string, responseContent) {
	token, err := jwt.Parse(tokenString, getTokenParsingHandler(dbName))
	if err != nil {
		return "", "", noToken(messageKey, dbName, request)
	} else {
		if token == nil {
			return "", "", notAuthorization(messageKey, dbName, request)
		} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userUuid, user, tokenExpired := getDataFromJWTClaims(claims)
			if tokenExpired {
				return "", "", expired(messageKey, dbName, user, request)
			} else {
				return userUuid, user, handleActions(dbName, userUuid, user, request)
			}
		} else {
			return "", "", invalidToken(messageKey, dbName, request)
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

func invalidMessage(request ApiParameters) responseContent {
	return responseContent{
		Status:      FailedStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		TextMessage: "Incorrect message",
	}
}

func handleActions(dbName string, userUuid string, userName string, request ApiParameters) responseContent {
	if request.Action == ActionLoad {
		return executeActionGrid(dbName, userUuid, userName, request, GetGridsRows)
	} else if request.Action == ActionChangeGrid {
		return executeActionGrid(dbName, userUuid, userName, request, PostGridsRows)
	} else if request.Action == ActionLocateGrid {
		return locate(request)
	} else {
		return invalidAction(dbName, request)
	}
}

func heartBeat(request ApiParameters) responseContent {
	return responseContent{
		Status: SuccessStatus,
		Action: request.Action,
	}
}

func logOut(request ApiParameters) responseContent {
	return responseContent{
		Status:      SuccessStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		TextMessage: "User logged out",
	}
}

func locate(request ApiParameters) responseContent {
	return responseContent{
		Status:     SuccessStatus,
		Action:     request.Action,
		ActionText: request.ActionText,
		GridUuid:   request.GridUuid,
		ColumnUuid: request.ColumnUuid,
		Uuid:       request.Uuid,
	}
}

func notAuthorization(messageKey, dbName string, request ApiParameters) responseContent {
	configuration.LogError(dbName, "", "No authorization for message %s action: %s %s", messageKey, request.Action, request.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		GridUuid:    request.GridUuid,
		TextMessage: "No authorization",
	}
}

func expired(messageKey, dbName string, userName string, request ApiParameters) responseContent {
	configuration.LogError(dbName, userName, "Authorization expired for message %s action: %s %s", messageKey, request.Action, request.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		GridUuid:    request.GridUuid,
		TextMessage: "Authorization expired",
	}
}

func invalidAction(dbName string, request ApiParameters) responseContent {
	configuration.Log(dbName, "", "Invalid action: %s %s", request.Action, request.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		GridUuid:    request.GridUuid,
		TextMessage: "Invalid action (" + request.Action + ")",
	}
}

func noToken(messageKey, dbName string, request ApiParameters) responseContent {
	configuration.LogError(dbName, "", "No token for message %s action: %s %s", messageKey, request.Action, request.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		GridUuid:    request.GridUuid,
		TextMessage: "Missing authorization",
	}
}

func invalidToken(messageKey, dbName string, request ApiParameters) responseContent {
	configuration.LogError(dbName, "", "Invalid token for message %s action: %s %s", messageKey, request.Action, request.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      request.Action,
		ActionText:  request.ActionText,
		GridUuid:    request.GridUuid,
		TextMessage: "Invalid token",
	}
}

type ActionGridDataFunc func(ct context.Context, p ApiParameters, payload GridPost) GridResponse

func executeActionGrid(dbName string, userUuid string, userName string, request ApiParameters, f ActionGridDataFunc) responseContent {
	request.DbName = dbName
	request.UserUuid = userUuid
	request.UserName = userName
	response := f(context.Background(), request, request.DataSet)
	if response.Err != nil {
		return responseContent{
			Status:      FailedStatus,
			Action:      request.Action,
			ActionText:  request.ActionText,
			GridUuid:    request.GridUuid,
			Uuid:        request.Uuid,
			TextMessage: response.Err.Error(),
		}
	}
	return responseContent{
		Status:     SuccessStatus,
		Action:     request.Action,
		ActionText: request.ActionText,
		GridUuid:   request.GridUuid,
		Uuid:       request.Uuid,
		DataSet:    response,
	}
}
