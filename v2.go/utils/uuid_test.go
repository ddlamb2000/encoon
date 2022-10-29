// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"testing"
)

func TestGetNewUUID(t *testing.T) {
	uuid := GetNewUUID()
	if uuid == "" {
		t.Fatal("No uuid generated.")
	}
}
