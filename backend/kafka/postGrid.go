// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"context"

	"d.lambert.fr/encoon/apis"
)

func postGrid(dbName string, userUuid string, userName string, content requestContent) responseContent {
	parameters := getParameters(dbName, userUuid, userName, content)
	response := apis.PostGridsRows(context.Background(), "", parameters, content.DataSet)
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
