// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"strings"
	"testing"
)

func TestSetPath(t *testing.T) {
	row := Row{
		Uuid: "12345",
	}
	row.SetPath("test", "users")
	expect := "/test/users/12345"
	if row.Path != expect {
		t.Errorf(`Got %v instead of %v.`, row.Path, expect)
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

func TestGetRowsColumnDefinitions(t *testing.T) {
	got := getRowsColumnDefinitions()
	expect := "text01 text, text02 text, text03 text"
	if !strings.Contains(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
	expect2 := "int01 integer, int02 integer, int03 integer"
	if !strings.Contains(got, expect2) {
		t.Errorf(`Got %v instead of %v.`, got, expect2)
	}
}
