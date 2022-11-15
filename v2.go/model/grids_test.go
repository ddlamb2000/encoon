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
