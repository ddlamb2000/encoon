// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

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
