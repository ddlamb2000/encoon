// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestGridSetPath(t *testing.T) {
	grid := GetNewGrid()
	grid.Uuid = "xxx"
	text1 := "aaa"
	grid.Text1 = &text1
	grid.SetPathAndDisplayString("test")
	expect := "/test/xxx"
	if grid.Path != expect {
		t.Errorf(`Got %v instead of %v.`, grid.Path, expect)
	}
	expect = "aaa"
	if grid.DisplayString != expect {
		t.Errorf(`Got %v instead of %v.`, grid.DisplayString, expect)
	}
}

func TestGetTableName(t *testing.T) {
	tests := []struct {
		uuid, expect string
	}{
		{"1234", "rows"},
		{UuidGrids, "grids"},
		{UuidColumns, "columns"},
		{UuidRelationships, "relationships"},
		{UuidMigrations, "migrations"},
		{UuidUsers, "users"},
	}
	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			grid := GetNewGrid()
			grid.Uuid = tt.uuid
			got := grid.GetTableName()
			if got != tt.expect {
				t.Errorf(`Got %v instead of %v.`, got, tt.expect)
			}
		})
	}
}

func TestSetViewEditAccessFlags(t *testing.T) {
	tests := []struct {
		test              string
		ownerUuid         string
		defaultAccessUuid string
		viewAccessUuid    string
		editAccessUuid    string
		userUuid          string
		expectCanView     bool
		expectCanEdit     bool
		expectGetCanView  bool
		expectGetCanEdit  bool
	}{
		{"1", "user1", "", "", "", "user1", true, true, true, true},
		{"2", "user1", "", "", "", "user2", false, false, false, false},
		{"3", "user1", UuidAccessLevelReadAccess, "", "", "user2", true, false, true, false},
		{"4", "user1", UuidAccessLevelWriteAccess, "", "", "user2", true, true, true, true},
		{"5", "user1", "", "user2", "", "user2", true, false, true, false},
		{"6", "user1", "", "", "user2", "user2", true, true, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid()
			grid.Owners[tt.ownerUuid] = true
			grid.DefaultAccess[tt.defaultAccessUuid] = true
			grid.ViewAccess[tt.viewAccessUuid] = true
			grid.EditAccess[tt.editAccessUuid] = true
			grid.SetViewEditAccessFlags(tt.userUuid)
			if grid.CanView != tt.expectCanView {
				t.Errorf(`Got grid.CanView=%v instead of %v.`, grid.CanView, tt.expectCanView)
			}
			if grid.CanEdit != tt.expectCanEdit {
				t.Errorf(`Got grid.CanEdit=%v instead of %v.`, grid.CanEdit, tt.expectCanEdit)
			}
			if grid.GetCanView() != tt.expectGetCanView {
				t.Errorf(`Got grid.GetCanView()=%v instead of %v.`, grid.GetCanView(), tt.expectCanView)
			}
			if grid.GetCanEdit() != tt.expectGetCanEdit {
				t.Errorf(`Got grid.GetCanEdit()=%v instead of %v.`, grid.GetCanEdit(), tt.expectCanEdit)
			}
		})
	}
}
