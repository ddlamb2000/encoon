// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

func locate(dbName string, content requestContent) responseContent {
	return responseContent{
		Status:     SuccessStatus,
		Action:     content.Action,
		GridUuid:   content.GridUuid,
		ColumnUuid: content.ColumnUuid,
		RowUuid:    content.RowUuid,
	}
}
