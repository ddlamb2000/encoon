// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestCleanupStrings(t *testing.T) {
	tests := []struct {
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
	got, _ := CalculateFileHash(fileName)
	expect := "468d2c16ea5eb8fe2f2244faa7c3307c21ec512c05a8b83ecb17c56c3812af4e"
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestCalculateFileHash2(t *testing.T) {
	fileName := "../testData/validConfiguration2.yml"
	got, _ := CalculateFileHash(fileName)
	expect := "548641131cdcb580f58aa24e1ece78fa4167234d0b45f73e9abf5420e351f28d"
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestCalculateFileHash3(t *testing.T) {
	fileName := "../xxx.yml"
	got, _ := CalculateFileHash(fileName)
	expect := ""
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}
