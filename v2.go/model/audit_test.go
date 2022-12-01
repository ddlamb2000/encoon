// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestGetNewAudit(t *testing.T) {
	audit := GetNewAudit()
	if audit == nil {
		t.Errorf(`Isse when creating audit.`)
	}
}

func TestActionName(t *testing.T) {
	tests := []struct {
		test             string
		columnName       string
		expectActionName string
	}{
		{"1", "relationship1", "Created"},
		{"2", "relationship2", "Updated"},
		{"3", "relationship3", "Deleted"},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			audit := GetNewAudit()
			audit.ColumnName = tt.columnName
			audit.SetActionName()
			if audit.ActionName != tt.expectActionName {
				t.Errorf(`Got audit.ActionName=%v instead of %v.`, audit.ActionName, tt.expectActionName)
			}
		})
	}
}
