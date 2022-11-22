// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestGridSetPath(t *testing.T) {
	grid := Grid{}
	grid.Uuid = "xxx"
	grid.SetPath("test")
	expect := "/test/xxx"
	if grid.Path != expect {
		t.Errorf(`Got %v instead of %v.`, grid.Path, expect)
	}
}
