// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"reflect"
	"testing"
)

func TestGetNewRow(t *testing.T) {
	row := GetNewRow()
	if row == nil {
		t.Errorf(`Issue when creating row.`)
	}
}

func TestGetNewRowWithUuid(t *testing.T) {
	row := GetNewRowWithUuid()
	if row == nil || row.Uuid == "" {
		t.Errorf(`Issue when creating row Uuid.`)
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

func TestRowSetViewEditAccessFlags(t *testing.T) {
	user1 := "aaaa"
	grid1 := GetNewGrid("")
	grid2 := GetNewGrid("")
	grid2.ViewAccess[user1] = true
	grid3 := GetNewGrid("")
	grid3.EditAccess[user1] = true
	grid4 := GetNewGrid("")
	grid4.Uuid = UuidColumns
	grid4.Owners[user1] = true
	grid5 := GetNewGrid(UuidGrids)
	grid6 := GetNewGrid(UuidUsers)
	tests := []struct {
		test             string
		grid             *Grid
		userUuid         string
		rowUuid          string
		expectCanViewRow bool
		expectCanEditRow bool
	}{
		{"1", grid1, "aaaa", "xxx", false, false},
		{"2", grid2, "aaaa", "xxx", true, false},
		{"3", grid3, "aaaa", "xxx", true, true},
		{"4", grid4, "aaaa", "xxx", true, true},
		{"5", grid4, "bbbb", "xxx", true, false},
		{"6", grid5, "bbbb", "xxx", true, false},
		{"7", grid6, "bbbb", "xxx", false, false},
		{"8", grid6, "xxx", "xxx", true, true},
		{"9", nil, "aaaa", "xxx", true, true},
		{"10", nil, "bbbb", "xxx", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			row := GetNewRow()
			row.Uuid = tt.rowUuid
			if tt.grid != nil {
				row.GridUuid = tt.grid.Uuid
			}
			row.CreatedBy = &user1
			row.SetViewEditAccessFlags(tt.grid, tt.userUuid)
			if row.CanViewRow != tt.expectCanViewRow {
				t.Errorf(`Got row.CanViewRow=%v instead of %v.`, row.CanViewRow, tt.expectCanViewRow)
			}
			if row.CanEditRow != tt.expectCanEditRow {
				t.Errorf(`Got row.CanEditRow=%v instead of %v.`, row.CanEditRow, tt.expectCanEditRow)
			}
		})
	}
}

func TestGetRowsQueryOutput(t *testing.T) {
	row := GetNewRow()
	tests := []struct {
		test   string
		uuid   string
		expect []any
	}{
		{"1", UuidGrids, []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Text2,
			&row.Text3,
			&row.Enabled,
			&row.Revision,
		}},
		{"2", UuidColumns, []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Text2,
			&row.Text3,
			&row.Int1,
			&row.Enabled,
			&row.Revision,
		}},
		{"3", UuidRelationships, []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Text2,
			&row.Text3,
			&row.Text4,
			&row.Text5,
			&row.Enabled,
			&row.Revision,
		}},
		{"4", UuidMigrations, []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Int1,
			&row.Enabled,
			&row.Revision,
		}},
		{"5", UuidUsers, []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Text2,
			&row.Text3,
			&row.Text4,
			&row.Enabled,
			&row.Revision,
		}},
		{"6", UuidTransactions, []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Enabled,
			&row.Revision,
		}},
		{"7", "xxx", []any{
			&row.Uuid,
			&row.GridUuid,
			&row.Created,
			&row.CreatedBy,
			&row.Updated,
			&row.UpdatedBy,
			&row.Text1,
			&row.Text2,
			&row.Text3,
			&row.Text4,
			&row.Text5,
			&row.Text6,
			&row.Text7,
			&row.Text8,
			&row.Text9,
			&row.Text10,
			&row.Int1,
			&row.Int2,
			&row.Int3,
			&row.Int4,
			&row.Int5,
			&row.Int6,
			&row.Int7,
			&row.Int8,
			&row.Int9,
			&row.Int10,
			&row.Enabled,
			&row.Revision,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			row.GridUuid = tt.uuid
			got := row.GetRowsQueryOutput()
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf(`Got %s instead of %s.`, got, tt.expect)
			}
		})
	}
}

func TestAppendRowValuesForSeedRowDb(t *testing.T) {
	row := GetNewRow()
	tests := []struct {
		test   string
		uuid   string
		expect []any
	}{
		{"1", UuidGrids, []any{
			row.Text1,
			row.Text2,
			row.Text3,
		}},
		{"2", UuidColumns, []any{
			row.Text1,
			row.Text2,
			row.Text3,
			row.Int1,
		}},
		{"3", UuidRelationships, []any{
			row.Text1,
			row.Text2,
			row.Text3,
			row.Text4,
			row.Text5,
		}},
		{"4", "xxx", []any{
			row.Text1,
			row.Text2,
			row.Text3,
			row.Text4,
			row.Text5,
			row.Text6,
			row.Text7,
			row.Text8,
			row.Text9,
			row.Text10,
			row.Int1,
			row.Int2,
			row.Int3,
			row.Int4,
			row.Int5,
			row.Int6,
			row.Int7,
			row.Int8,
			row.Int9,
			row.Int10,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			values := make([]any, 0)
			row.GridUuid = tt.uuid
			got := row.AppendRowValuesForSeedRowDb(values)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf(`Got %s instead of %s.`, got, tt.expect)
			}
		})
	}
}

func TestRowGetValueAsString(t *testing.T) {
	tests := []struct {
		test       string
		columnName string
		expect     string
	}{
		{"1", "text1", "A"},
		{"2", "text2", "B"},
		{"3", "text3", "C"},
		{"4", "text4", "A"},
		{"5", "text5", "B"},
		{"6", "text6", "C"},
		{"7", "text7", "A"},
		{"8", "text8", "B"},
		{"9", "text9", "C"},
		{"10", "text10", "A"},
		{"11", "int1", "1"},
		{"12", "int2", "2"},
		{"13", "int3", "3"},
		{"14", "int4", "1"},
		{"15", "int5", "2"},
		{"16", "int6", "3"},
		{"17", "int7", "1"},
		{"18", "int8", "2"},
		{"19", "int9", "3"},
		{"20", "int10", "1"},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			text1 := "A"
			text2 := "B"
			text3 := "C"
			i1 := int64(1)
			i2 := int64(2)
			i3 := int64(3)
			row := Row{
				Uuid:     "12345",
				GridUuid: "56789",
				Text1:    &text1,
				Text2:    &text2,
				Text3:    &text3,
				Text4:    &text1,
				Text5:    &text2,
				Text6:    &text3,
				Text7:    &text1,
				Text8:    &text2,
				Text9:    &text3,
				Text10:   &text1,
				Int1:     &i1,
				Int2:     &i2,
				Int3:     &i3,
				Int4:     &i1,
				Int5:     &i2,
				Int6:     &i3,
				Int7:     &i1,
				Int8:     &i2,
				Int9:     &i3,
				Int10:    &i1,
				Enabled:  false,
			}
			got := row.GetValueAsString(tt.columnName)
			if got != tt.expect {
				t.Errorf(`Got row.GetValueAsString=%v instead of %v.`, got, tt.expect)
			}
		})
	}
}
