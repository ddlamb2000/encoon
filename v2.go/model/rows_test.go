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

func TestRowSetViewEditAccessFlags(t *testing.T) {
	user1 := "aaaa"
	grid1 := GetNewGrid()
	grid2 := GetNewGrid()
	grid2.CanViewRows = true
	grid3 := GetNewGrid()
	grid3.CanEditRows = true
	grid4 := GetNewGrid()
	grid4.CanEditOwnedRows = true
	grid4.Owners["aaaa"] = true
	tests := []struct {
		test             string
		grid             *Grid
		uuid             string
		expectCanViewRow bool
		expectCanEditRow bool
	}{
		{"1", grid1, "aaaa", false, false},
		{"2", grid2, "aaaa", true, false},
		{"3", grid3, "aaaa", false, true},
		{"4", grid4, "aaaa", false, true},
		{"5", grid4, "bbbb", false, false},
		{"6", nil, "aaaa", true, true},
		{"7", nil, "bbbb", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			row := GetNewRow()
			row.CreatedBy = &user1
			row.SetViewEditAccessFlags(tt.grid, tt.uuid)
			if row.CanViewRow != tt.expectCanViewRow {
				t.Errorf(`Got row.CanViewRowsRow=%v instead of %v.`, row.CanViewRow, tt.expectCanViewRow)
			}
			if row.CanEditRow != tt.expectCanEditRow {
				t.Errorf(`Got row.CanEditRowsRow=%v instead of %v.`, row.CanEditRow, tt.expectCanEditRow)
			}
		})
	}
}
