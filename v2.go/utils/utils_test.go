// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	Log("Test: %v", "test")
}

func TestLogError(t *testing.T) {
	LogError("Test: %v", "test")
}

func TestLogAndReturnError(t *testing.T) {
	got := LogAndReturnError("Test: %v", "test")
	expect := "Test: test"
	if got.Error() != expect {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestTrace(t *testing.T) {
	Trace("true", "Test: %v", "test")
}

func TestCleanupStrings(t *testing.T) {
	var tests = []struct {
		id     int
		given  string
		expect string
	}{
		{1, `{           "message":       "Not authorized."}`, `{ "message": "Not authorized."}`},
		{2, `{"message": "Not authorized."}`, `{"message": "Not authorized."}`},
		{3, `{      "message": "Not authorized."     }`, `{ "message": "Not authorized." }`},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.id)
		t.Run(testname, func(t *testing.T) {
			then := CleanupStrings(tt.given)
			if then != tt.expect {
				t.Errorf("Got %q instead of %q.", then, tt.expect)
			}
		})
	}
}

func TestCalculateFileHash1(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	got := CalculateFileHash(fileName)
	expect := "b65c60887f0773546cdb387a71fcf9ce2b8c49008c6a4ebb6923f0dafb585253"
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestCalculateFileHash2(t *testing.T) {
	fileName := "../testData/validConfiguration2.yml"
	got := CalculateFileHash(fileName)
	expect := "548641131cdcb580f58aa24e1ece78fa4167234d0b45f73e9abf5420e351f28d"
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}
