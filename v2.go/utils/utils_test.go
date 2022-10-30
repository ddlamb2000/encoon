// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"testing"
)

func TestInitWithLog(t *testing.T) {
	InitWithLog()
}

func TestCleanupStrings(t *testing.T) {
	var tests = []struct {
		id       int
		given    string
		expected string
	}{
		{1, `{           "message":       "Not authorized."}`, `{ "message": "Not authorized."}`},
		{2, `{"message": "Not authorized."}`, `{"message": "Not authorized."}`},
		{3, `{      "message": "Not authorized."     }`, `{ "message": "Not authorized." }`},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.id)
		t.Run(testname, func(t *testing.T) {
			then := CleanupStrings(tt.given)
			if then != tt.expected {
				t.Errorf("Got %q instead of %q.", then, tt.expected)
			}
		})
	}
}
