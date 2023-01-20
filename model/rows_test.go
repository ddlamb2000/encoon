// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"testing"
)

func TestGetNewRow(t *testing.T) {
	row := GetNewRow()
	if row == nil {
		t.Errorf(`Isseu when creating row.`)
	}
}

func TestDisplayString(t *testing.T) {
	text1 := "xxx"
	row := Row{
		Uuid:     "12345",
		GridUuid: "56789",
		Text1:    &text1,
		Enabled:  true,
	}
	row.SetDisplayString("test")
	expect := "xxx"
	if row.DisplayString != expect {
		t.Errorf(`Got %v instead of %v.`, row.DisplayString, expect)
	}
}

func TestDisplayString2(t *testing.T) {
	text1 := "xxx"
	row := Row{
		Uuid:     "12345",
		GridUuid: "56789",
		Text1:    &text1,
		Enabled:  false,
	}
	row.SetDisplayString("test")
	expect := "xxx [DELETED]"
	if row.DisplayString != expect {
		t.Errorf(`Got %v instead of %v.`, row.DisplayString, expect)
	}
}

func TestDisplayString3(t *testing.T) {
	row := Row{
		Uuid:     "12345",
		GridUuid: "56789",
		Enabled:  true,
	}
	row.SetDisplayString("test")
	expect := "12345"
	if row.DisplayString != expect {
		t.Errorf(`Got %v instead of %v.`, row.DisplayString, expect)
	}
}

func TestDisplayString5(t *testing.T) {
	var i1 int64 = 10
	row := Row{
		Uuid:     "12345",
		GridUuid: "56789",
		Int1:     &i1,
		Enabled:  true,
	}
	row.SetDisplayString("test")
	expect := "10"
	if row.DisplayString != expect {
		t.Errorf(`Got %v instead of %v.`, row.DisplayString, expect)
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

func TestRowSetViewEditAccessFlags(t *testing.T) {
	user1 := "aaaa"
	grid1 := GetNewGrid()
	grid2 := GetNewGrid()
	grid2.ViewAccess[user1] = true
	grid3 := GetNewGrid()
	grid3.EditAccess[user1] = true
	grid4 := GetNewGrid()
	grid4.Uuid = UuidColumns
	grid4.Owners[user1] = true
	tests := []struct {
		test             string
		grid             *Grid
		uuid             string
		expectCanViewRow bool
		expectCanEditRow bool
	}{
		{"1", grid1, "aaaa", false, false},
		{"2", grid2, "aaaa", true, false},
		{"3", grid3, "aaaa", true, false},
		{"4", grid4, "aaaa", true, true},
		{"5", grid4, "bbbb", true, false},
		{"6", nil, "aaaa", true, true},
		{"7", nil, "bbbb", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			row := GetNewRowWithUuid()
			row.CreatedBy = &user1
			row.SetViewEditAccessFlags(tt.grid, tt.uuid)
			if row.CanViewRow != tt.expectCanViewRow {
				t.Errorf(`Got row.CanViewRow=%v instead of %v.`, row.CanViewRow, tt.expectCanViewRow)
			}
			if row.CanEditRow != tt.expectCanEditRow {
				t.Errorf(`Got row.CanEditRow=%v instead of %v.`, row.CanEditRow, tt.expectCanEditRow)
			}
		})
	}
}
