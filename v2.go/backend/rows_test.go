// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import "testing"

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
