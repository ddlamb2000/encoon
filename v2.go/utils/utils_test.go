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
	expect := "984795033e28dab99b3cbb0ff6e8f96e0c514599cd773157da92ffaf52c2b859"
	if !strings.Contains(got, expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestCalculateFileHash2(t *testing.T) {
	fileName := "../testData/validConfiguration2.yml"
	got, _ := CalculateFileHash(fileName)
	expect := "a2dde6ccb1b80b8f461641b9cc92a5f50097626334c6f8759653b9c8b8930dff"
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
