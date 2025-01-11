// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestGetDataFromHeaders(t *testing.T) {
	message := kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "requestInitiatedOn",
				Value: []byte("A"),
			},
			{
				Key:   "jwt",
				Value: []byte("B"),
			},
			{
				Key:   "gridUuid",
				Value: []byte("C"),
			},
			{
				Key:   "contextUuid",
				Value: []byte("D"),
			},
		},
	}

	requestInitiatedOn, tokenString, gridUuid, contextUuid := getDataFromHeaders(message)
	expectRequestInitiatedOn, expectTokenString, expectGridUuid, expectContextUuid := "A", "B", "C", "D"

	if requestInitiatedOn != expectRequestInitiatedOn {
		t.Errorf(`Got %s instead of %s.`, requestInitiatedOn, expectRequestInitiatedOn)
	}
	if tokenString != expectTokenString {
		t.Errorf(`Got %s instead of %s.`, tokenString, expectTokenString)
	}
	if gridUuid != expectGridUuid {
		t.Errorf(`Got %s instead of %s.`, gridUuid, expectGridUuid)
	}
	if contextUuid != expectContextUuid {
		t.Errorf(`Got %s instead of %s.`, contextUuid, expectContextUuid)
	}
}
