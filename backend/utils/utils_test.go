// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

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
	expect := "d335029e389481c488c28b9ac406fe3e38cb8f04edddde6e1383cb8a5c31f1b6"
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestCalculateFileHash2(t *testing.T) {
	fileName := "../testData/validConfiguration2.yml"
	got, _ := CalculateFileHash(fileName)
	expect := "5c6baaaa0393ffa7584c8b0318908ea16c792552bd4027d36c893ef765546dac"
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
