// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"testing"
)

func TestGetNewColumn(t *testing.T) {
	column := GetNewColumn()
	if column == nil {
		t.Errorf(`Isse when creating column.`)
	}
}

func TestColumnString(t *testing.T) {
	column := Column{Name: "text1", Label: "Label"}
	got := column.String()
	expect := `text1 "Label"`
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestIsAttribute(t *testing.T) {
	column := Column{TypeUuid: UuidTextColumnType}
	got := column.IsAttribute()
	expect := true
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestIsAttribute2(t *testing.T) {
	column := Column{TypeUuid: UuidReferenceColumnType}
	got := column.IsAttribute()
	expect := false
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestIsReference(t *testing.T) {
	column := Column{TypeUuid: UuidReferenceColumnType}
	got := column.IsReference()
	expect := true
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestIsReference2(t *testing.T) {
	column := Column{TypeUuid: UuidTextColumnType}
	got := column.IsReference()
	expect := false
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestIsOwned(t *testing.T) {
	column := Column{Owned: true}
	got := column.IsOwned()
	expect := true
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetColumnNamePrefixFromType(t *testing.T) {
	tests := []struct {
		test       string
		columnType string
		expect     string
	}{
		{"1", UuidBooleanColumnType, "text"},
		{"2", UuidIntColumnType, "int"},
		{"3", UuidPasswordColumnType, "text"},
		{"4", UuidReferenceColumnType, "relationship"},
		{"5", UuidRichTextColumnType, "text"},
		{"6", UuidTextColumnType, "text"},
		{"7", UuidUuidColumnType, "text"},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			column := GetNewColumn()
			column.TypeUuid = tt.columnType
			got := column.GetColumnNamePrefixFromType()
			if got != tt.expect {
				t.Errorf(`Got %v instead of %v.`, got, tt.expect)
			}
		})
	}
}

func TestGetColumnNamePrefixAndIndex(t *testing.T) {
	tests := []struct {
		test         string
		name         string
		expectPrefix string
		expectIndex  int64
	}{
		{"1", "text1", "text", 1},
		{"2", "text2", "text", 2},
		{"3", "text123", "text", 123},
		{"4", "text", "", 0},
		{"5", "int1", "int", 1},
		{"6", "int24", "int", 24},
		{"7", "int", "", 0},
		{"8", "relationship1", "relationship", 1},
		{"9", "relationship54", "relationship", 54},
		{"10", "relationship", "", 0},
		{"11", "xxx2", "", 0},
		{"12", "y", "", 0},
		{"13", "", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			column := GetNewColumn()
			column.Name = tt.name
			gotPrefix, gotIndex := column.GetColumnNamePrefixAndIndex()
			if gotPrefix != tt.expectPrefix {
				t.Errorf(`Got prefix %s instead of %s from %s.`, gotPrefix, tt.expectPrefix, tt.name)
			}
			if gotIndex != tt.expectIndex {
				t.Errorf(`Got index %d instead of %d from %s.`, gotIndex, tt.expectIndex, tt.name)
			}
		})
	}
}
