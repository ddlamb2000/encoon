// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"testing"
)

func TestGetNewGrid(t *testing.T) {
	grid := GetNewGrid()
	if grid == nil {
		t.Errorf(`Isse when creating grid.`)
	}
}

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

func TestGridSetViewEditAccessFlags(t *testing.T) {
	tests := []struct {
		test                   string
		uuid                   string
		ownerUuid              string
		defaultAccessUuid      string
		viewAccessUuid         string
		editAccessUuid         string
		userUuid               string
		expectCanViewRows      bool
		expectCanEditRows      bool
		expectCanEditOwnedRows bool
		expectCanAddRows       bool
	}{
		{"1", "aaaa", "user1", "", "", "", "user1", true, true, true, true},
		{"2", "aaaa", "user1", "", "", "", "user2", false, false, false, false},
		{"3", "aaaa", "user1", UuidAccessLevelReadAccess, "", "", "user2", true, false, false, false},
		{"4", "aaaa", "user1", UuidAccessLevelWriteAccess, "", "", "user2", true, true, true, true},
		{"5", "aaaa", "user1", "", "user2", "", "user2", true, false, false, false},
		{"6", "aaaa", "user1", "", "", "user2", "user2", true, true, true, true},
		{"7", UuidGrids, "", "", "", "", "user1", true, false, true, true},
		{"8", UuidColumns, "", "", "", "", "user1", true, false, true, true},
		{"9", UuidUsers, "", "", "", "", "user1", true, false, false, false},
		{"10", UuidAccessLevel, "", "", "", "", "user1", true, false, false, false},
		{"11", UuidColumnTypes, "", "", "", "", "user1", true, false, false, false},
		{"12", UuidMigrations, "", "", "", "", "user1", false, false, false, false},
		{"13", UuidRelationships, "", "", "", "", "user1", true, false, true, true},
		{"14", UuidTransactions, "", "", "", "", "user1", false, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid()
			grid.Uuid = tt.uuid
			grid.Owners[tt.ownerUuid] = true
			grid.DefaultAccess[tt.defaultAccessUuid] = true
			grid.ViewAccess[tt.viewAccessUuid] = true
			grid.EditAccess[tt.editAccessUuid] = true
			grid.SetViewEditAccessFlags(tt.userUuid)
			if grid.CanViewRows != tt.expectCanViewRows {
				t.Errorf(`Got grid.CanViewRowsRows=%v instead of %v.`, grid.CanViewRows, tt.expectCanViewRows)
			}
			if grid.CanEditRows != tt.expectCanEditRows {
				t.Errorf(`Got grid.CanEditRowsRows=%v instead of %v.`, grid.CanEditRows, tt.expectCanEditRows)
			}
			if grid.CanEditOwnedRows != tt.expectCanEditOwnedRows {
				t.Errorf(`Got grid.CanEditOwnedRows=%v instead of %v.`, grid.CanEditOwnedRows, tt.expectCanEditOwnedRows)
			}
			if grid.CanAddRows != tt.expectCanAddRows {
				t.Errorf(`Got grid.CanAddRowsRows=%v instead of %v.`, grid.CanAddRows, tt.expectCanAddRows)
			}
		})
	}
}

func TestCopyAccessToOtherGrid(t *testing.T) {
	uuid1, uuid2, uuid3, uuid4 := "aaa", "bbb", "ccc", "ddd"
	grid1 := GetNewGrid()
	grid1.OwnerUuid = &uuid1
	grid1.DefaultAccessUuid = &uuid2
	grid1.ViewAccessUuid = &uuid3
	grid1.EditAccessUuid = &uuid4
	grid2 := GetNewGrid()
	grid1.CopyAccessToOtherGrid(grid2)
	if !grid2.Owners[uuid1] {
		t.Errorf(`Can't find owner.`)
	}
	if !grid2.DefaultAccess[uuid2] {
		t.Errorf(`Can't find default access.`)
	}
	if !grid2.ViewAccess[uuid3] {
		t.Errorf(`Can't find view access.`)
	}
	if !grid2.EditAccess[uuid4] {
		t.Errorf(`Can't find edit access.`)
	}
}

func TestHasOwnership(t *testing.T) {
	grid := GetNewGrid()
	grid.Owners["aaaa"] = true
	grid.Owners["bbbb"] = true
	tests := []struct {
		test   string
		uuid   string
		expect bool
	}{
		{"1", "aaaa", true},
		{"2", "bbbb", true},
		{"3", "cccc", false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			if grid.HasOwnership(tt.uuid) != tt.expect {
				t.Errorf(`Got grid.HasOwnership=%v instead of %v.`, grid.HasOwnership(tt.uuid), tt.expect)
			}
		})
	}
}
