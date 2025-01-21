// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"testing"

	"d.lambert.fr/encoon/model"
)

func RunSystemTestKafka(t *testing.T) {
	t.Run("Heartbeat", func(t *testing.T) {
		response, _ := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, "", requestContent{
			Action: ActionHeartbeat,
		})
		responseIsSuccess(t, response)
	})

	t.Run("LoadGrid", func(t *testing.T) {
		response, responseData := runKafkaTestRequest(t, "test", "root", model.UuidRootUser, model.UuidGrids, requestContent{
			Action:   ActionLoad,
			GridUuid: model.UuidGrids,
		})
		responseIsSuccess(t, response)
		jsonStringContains(t, responseData, `"grid":{"gridUuid":"`+model.UuidGrids+`","uuid":"`+model.UuidGrids+`"`)
		jsonStringContains(t, responseData, `"rows":[`)
		jsonStringContains(t, responseData, `"createdBy":"`+model.UuidRootUser+`"`)
		jsonStringContains(t, responseData, `"columns":[`)
		jsonStringContains(t, responseData, `"label":"Columns","name":"relationship1","type":"Reference"`)
		jsonStringContains(t, responseData, `"label":"Description","name":"text2","type":"Text"`)
	})
}
