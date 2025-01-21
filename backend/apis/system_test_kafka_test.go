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
}
