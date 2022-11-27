// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestGetNewRow(t *testing.T) {
	row := GetNewRow()
	if row == nil {
		t.Errorf(`Isse when creating row.`)
	}
}

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
		Uuid:          "12345",
		DisplayString: "xyz",
	}
	got := row.String()
	expect := "xyz [12345]"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetRowsQueryOutput(t *testing.T) {
	row := Row{
		Uuid:          "12345",
		DisplayString: "xyz",
	}
	got := len(row.GetRowsQueryOutput())
	expect := 28
	if got != expect {
		t.Errorf(`Got %d instead of %d.`, got, expect)
	}
}
