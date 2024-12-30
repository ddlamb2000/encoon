// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"

	"d.lambert.fr/encoon/apis"
)

func getGrid(dbName string, userUuid string, userName string, content requestContent) responseContent {
	parameters := getParameters(dbName, userUuid, userName, content)
	response := apis.GetGridsRows(context.Background(), "", parameters)
	if response.Err != nil {
		return responseContent{
			Status:      FailedStatus,
			Action:      content.Action,
			TextMessage: response.Err.Error(),
		}
	}
	return responseContent{
		Status:   SuccessStatus,
		Action:   content.Action,
		GridUuid: content.GridUuid,
		DataSet:  response,
	}
}

func getParameters(dbName string, userUuid string, userName string, content requestContent) apis.ApiParameters {
	gridUuid, uuid := content.GridUuid, content.Uuid
	return apis.ApiParameters{
		DbName:   dbName,
		UserUuid: userUuid,
		UserName: userName,
		GridUuid: gridUuid,
		Uuid:     uuid,
	}
}
