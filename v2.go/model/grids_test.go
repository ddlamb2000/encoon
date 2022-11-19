// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestGridSetPath(t *testing.T) {
	grid := Grid{}
	grid.SetPath("test", "users")
	expect := "/test/users"
	if grid.Path != expect {
		t.Errorf(`Got %v instead of %v.`, grid.Path, expect)
	}
}

func TestGetUri(t *testing.T) {
	text1 := "test"
	grid := Grid{}
	grid.Text1 = &text1
	expect := "test"
	got := grid.GetUri()
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUri2(t *testing.T) {
	grid := Grid{}
	expect := "?"
	got := grid.GetUri()
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
