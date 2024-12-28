// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"

	"d.lambert.fr/encoon/apis"
)

func getGrid(dbName string, content requestContent) responseContent {
	p := getParameters(dbName, content)
	response := apis.GetGridsRows(context.Background(), "", p)
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

func getParameters(dbName string, content requestContent) apis.HtmlParameters {
	gridUuid, uuid := content.GridUuid, content.Uuid
	return apis.HtmlParameters{
		DbName:   dbName,
		GridUuid: gridUuid,
		Uuid:     uuid,
	}
}
