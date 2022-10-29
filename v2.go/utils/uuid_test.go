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

func TestCleanupStrings(t *testing.T) {
	given := `{           "message":       "Not authorized."}`
	then := CleanupStrings(given)
	expected := `{ "message": "Not authorized."}`
	if then != expected {
		t.Fatalf("Got %q instead of %q.", then, expected)
	}
}
