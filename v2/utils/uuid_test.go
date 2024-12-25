// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package utils

import (
	"testing"
)

func TestGetNewUUID(t *testing.T) {
	uuid := GetNewUUID()
	if uuid == "" {
		t.Errorf("No uuid generated.")
	}
	got := len(uuid)
	expect := 36
	if got != expect {
		t.Errorf("Uuid has a length of %d instead of %d.", got, expect)
	}
}
