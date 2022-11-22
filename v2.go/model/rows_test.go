// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestSetPath(t *testing.T) {
	text1 := "xxx"
	row := Row{
		Uuid:     "12345",
		GridUuid: "56789",
		Text1:    &text1,
	}
	row.SetPathAndDisplayString("test")
	expect := "/test/56789/12345"
	if row.Path != expect {
		t.Errorf(`Got %v instead of %v.`, row.Path, expect)
	}
	expect2 := "xxx"
	if row.DisplayString != expect2 {
		t.Errorf(`Got %v instead of %v.`, row.DisplayString, expect2)
	}
}

func TestRowAsString(t *testing.T) {
	row := Row{
		Uuid: "12345",
	}
	got := row.String()
	expect := "12345"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
