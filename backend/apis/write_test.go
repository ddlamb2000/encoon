// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestCalcHash(t *testing.T) {
	tests := []struct {
		test          string
		gridUuid      string
		expect        uint32
		expectBalance int
	}{
		{"1", "A", 3289118412, 0},
		{"2", "B", 3339451269, 0},
		{"3", "C", 3322673650, 1},
		{"4", "D", 3238785555, 0},
		{"5", "E", 3222007936, 1},
		{"6", "F", 3272340793, 1},
		{"7", "G", 3255563174, 2},
		{"8", "H", 3440116983, 0},
		{"9", "I", 3423339364, 1},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			got := calcHash(tt.gridUuid)
			if got != tt.expect {
				t.Errorf(`Got %d instead of %d.`, got, tt.expect)
			}
			nbPartitions := uint32(3)
			balance := int(got % nbPartitions)
			if balance != tt.expectBalance {
				t.Errorf(`Got %d instead of %d.`, balance, tt.expectBalance)
			}
		})
	}
}

func TestRoundRobinBalance(t *testing.T) {
	customRoundRobin := CustomRoundRobin{}

	tests := []struct {
		test     string
		gridUuid string
		expect   int
	}{
		{"1", "A", 0},
		{"2", "B", 0},
		{"3", "C", 1},
		{"4", "D", 0},
		{"5", "E", 1},
		{"6", "F", 1},
		{"7", "G", 2},
		{"8", "H", 0},
		{"9", "I", 1},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			message := kafka.Message{
				Headers: []kafka.Header{
					{
						Key:   "gridUuid",
						Value: []byte(tt.gridUuid),
					},
				},
			}
			got := customRoundRobin.Balance(message, 0, 1, 2)
			if got != tt.expect {
				t.Errorf(`Got %d instead of %d.`, got, tt.expect)
			}
		})
	}
}
