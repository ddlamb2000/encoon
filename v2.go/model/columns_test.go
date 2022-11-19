// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestColumnString(t *testing.T) {
	column := Column{Name: "text1", Label: "Label"}
	got := column.String()
	expect := `text1 "Label"`
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestIsAttribute(t *testing.T) {
	column := Column{TypeUuid: UuidColumns}
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
