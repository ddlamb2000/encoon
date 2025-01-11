// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

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

func ReadMessagesFromKafka() {
	kafkaBrokers := configuration.GetConfiguration().Kafka.Brokers
	for _, dbConfig := range configuration.GetConfiguration().Databases {
		go readMessages(dbConfig.Name, kafkaBrokers)
	}
}

func readMessages(dbName string, kafkaBrokers string) {
	topic := configuration.GetConfiguration().Kafka.TopicPrefix + "-" + dbName + "-requests"
	groupID := configuration.GetConfiguration().Kafka.GroupID + "-" + dbName
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
			configuration.LogError(dbName, "", "Error reading message: %v", err)
			continue
		}
		handleMessage(dbName, message)
	}
}

func handleMessage(dbName string, message kafka.Message) {
	requestReceivedOn := time.Now().UTC().Format(time.RFC3339Nano)
	requestInitiatedOn, tokenString, gridUuid, contextUuid := getDataFromHeaders(message)
	var content requestContent
	var response responseContent
	user := ""
	userUuid := ""
	if err := json.Unmarshal(message.Value, &content); err != nil {
		configuration.LogError(dbName, "", "Error message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText, err)
		response = invalidMessage(content)
	} else {
		configuration.Log(dbName, "", "PULL Message (%d bytes), topic: %s, key: %s, action: %s %s", len(message.Value), message.Topic, message.Key, content.Action, content.ActionText)
		if content.Action == ActionHeartbeat {
			response = heartBeat(content)
		} else if content.Action == ActionAuthentication {
			response = handleAuthentication(dbName, content)
		} else if content.Action == ActionLogout {
			response = logOut(content)
		} else {
			userUuid, user, response = validMessage(dbName, tokenString, content)
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

func validMessage(dbName string, tokenString string, content requestContent) (string, string, responseContent) {
	token, err := jwt.Parse(tokenString, getTokenParsingHandler(dbName))
	if err != nil {
		return "", "", invalidToken(dbName, content)
	} else {
		if token == nil {
			return "", "", notAuthorization(dbName, content)
		} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userUuid, user, tokenExpired := getDataFromJWTClaims(claims)
			if tokenExpired {
				return "", "", expired(dbName, user, content)
			} else {
				return userUuid, user, handleActions(dbName, userUuid, user, content)
			}
		} else {
			return "", "", invalidToken(dbName, content)
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
		RowUuid:    content.RowUuid,
	}
}

func notAuthorization(dbName string, content requestContent) responseContent {
	configuration.LogError(dbName, "", "No authorization for action: %s %s", content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		TextMessage: "No authorization",
	}
}

func expired(dbName string, userName string, content requestContent) responseContent {
	configuration.LogError(dbName, userName, "Authorization expired for action: %s %s", content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		TextMessage: "Authorization expired",
	}
}

func invalidAction(dbName string, content requestContent) responseContent {
	configuration.Log(dbName, "", "Invalid action: %s %s", content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
		TextMessage: "Invalid action (" + content.Action + ")",
	}
}

func invalidToken(dbName string, content requestContent) responseContent {
	configuration.LogError(dbName, "", "Invalid token for action: %s %s", content.Action, content.ActionText)
	return responseContent{
		Status:      FailedStatus,
		Action:      content.Action,
		ActionText:  content.ActionText,
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
			TextMessage: response.Err.Error(),
		}
	}
	return responseContent{
		Status:     SuccessStatus,
		Action:     content.Action,
		ActionText: content.ActionText,
		GridUuid:   content.GridUuid,
		DataSet:    response,
	}
}
