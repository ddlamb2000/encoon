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

func TestGridDisplayString(t *testing.T) {
	grid := GetNewGrid()
	grid.Uuid = "xxx"
	text1 := "aaa"
	grid.Text1 = &text1
	grid.SetDisplayString("test")
	expect := "aaa"
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
		{UuidTransactions, "transactions"},
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
		{"4", "aaaa", "user1", UuidAccessLevelWriteAccess, "", "", "user2", true, false, true, true},
		{"5", "aaaa", "user1", "", "user2", "", "user2", true, false, false, false},
		{"6", "aaaa", "user1", "", "", "user2", "user2", true, false, true, true},
		{"7", UuidGrids, "root", "", "", "", "user1", true, false, true, true},
		{"8", UuidColumns, "root", "", "", "", "user1", true, false, true, true},
		{"9", UuidUsers, "root", "", "", "", "user1", true, false, false, false},
		{"10", UuidAccessLevels, "root", "", "", "", "user1", true, false, false, false},
		{"11", UuidColumnTypes, "root", "", "", "", "user1", true, false, false, false},
		{"12", UuidMigrations, "root", "", "", "", "user1", false, false, false, false},
		{"13", UuidRelationships, "root", "", "", "", "user1", true, false, true, true},
		{"14", UuidTransactions, "root", "", "", "", "user1", false, false, false, false},
		{"15", UuidTransactions, "root", "", "", "", "root", true, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid()
			grid.Uuid = tt.uuid
			grid.Owners[tt.ownerUuid] = true
			grid.DefaultAccess[tt.defaultAccessUuid] = true
			grid.ViewAccess[tt.viewAccessUuid] = true
			grid.EditAccess[tt.editAccessUuid] = true
			canViewRows, canEditRows, canEditOwnedRows, canAddRows := grid.GetViewEditAccessFlags(tt.userUuid)
			if canViewRows != tt.expectCanViewRows {
				t.Errorf(`Got canViewRows=%v instead of %v.`, canViewRows, tt.expectCanViewRows)
			}
			if canEditRows != tt.expectCanEditRows {
				t.Errorf(`Got canEditRows=%v instead of %v.`, canEditRows, tt.expectCanEditRows)
			}
			if canEditOwnedRows != tt.expectCanEditOwnedRows {
				t.Errorf(`Got canEditOwnedRows=%v instead of %v.`, canEditOwnedRows, tt.expectCanEditOwnedRows)
			}
			if canAddRows != tt.expectCanAddRows {
				t.Errorf(`Got canAddRows=%v instead of %v.`, canAddRows, tt.expectCanAddRows)
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
