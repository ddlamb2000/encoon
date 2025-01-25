// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"reflect"
	"testing"
)

func TestGetNewGrid(t *testing.T) {
	grid := GetNewGrid("")
	if grid == nil {
		t.Errorf(`Isse when creating grid.`)
	}
}

func TestGridDisplayString(t *testing.T) {
	grid := GetNewGrid("xxx")
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
			grid := GetNewGrid(tt.uuid)
			got := grid.GetTableName()
			if got != tt.expect {
				t.Errorf(`Got %v instead of %v.`, got, tt.expect)
			}
		})
	}
}

func TestGridSetViewEditAccessFlags(t *testing.T) {
	tests := []struct {
		test              string
		uuid              string
		ownerUuid         string
		defaultAccessUuid string
		viewAccessUuid    string
		editAccessUuid    string
		userUuid          string
		expectCanViewRows bool
		expectCanEditRows bool
		expectCanAddRows  bool
		expectCanEditGrid bool
	}{
		{"1", "aaaa", "user1", "", "", "", "user1", true, true, true, true},
		{"2", "aaaa", "user1", "", "", "", "user2", false, false, false, false},
		{"3", "aaaa", "user1", UuidAccessLevelReadAccess, "", "", "user2", true, false, false, false},
		{"4", "aaaa", "user1", UuidAccessLevelWriteAccess, "", "", "user2", true, true, true, false},
		{"5", "aaaa", "user1", "", "user2", "", "user2", true, false, false, false},
		{"6", "aaaa", "user1", "", "", "user2", "user2", true, true, true, false},
		{"7", UuidGrids, "root", "", "", "", "user1", true, false, true, false},
		{"8", UuidColumns, "root", "", "", "", "user1", true, false, true, false},
		{"9", UuidUsers, "root", "", "", "", "user1", true, false, false, false},
		{"10", UuidAccessLevels, "root", "", "", "", "user1", true, false, false, false},
		{"11", UuidColumnTypes, "root", "", "", "", "user1", true, false, false, false},
		{"12", UuidMigrations, "root", "", "", "", "user1", false, false, false, false},
		{"13", UuidRelationships, "root", "", "", "", "user1", true, false, true, false},
		{"14", UuidTransactions, "root", "", "", "", "user1", false, false, false, false},
		{"15", UuidTransactions, "root", "", "", "", "root", true, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid(tt.uuid)
			grid.Owners[tt.ownerUuid] = true
			grid.DefaultAccess[tt.defaultAccessUuid] = true
			grid.ViewAccess[tt.viewAccessUuid] = true
			grid.EditAccess[tt.editAccessUuid] = true
			canViewRows, canEditRows, canAddRows, canEditGrid := grid.GetViewEditAccessFlags(tt.userUuid)
			if canViewRows != tt.expectCanViewRows {
				t.Errorf(`Got canViewRows=%v instead of %v.`, canViewRows, tt.expectCanViewRows)
			}
			if canEditRows != tt.expectCanEditRows {
				t.Errorf(`Got canEditRows=%v instead of %v.`, canEditRows, tt.expectCanEditRows)
			}
			if canAddRows != tt.expectCanAddRows {
				t.Errorf(`Got canAddRows=%v instead of %v.`, canAddRows, tt.expectCanAddRows)
			}
			if canEditGrid != tt.expectCanEditGrid {
				t.Errorf(`Got canEditGrid=%v instead of %v.`, canEditGrid, tt.expectCanEditGrid)
			}
		})
	}
}

func TestCopyAccessToOtherGrid(t *testing.T) {
	uuid1, uuid2, uuid3, uuid4 := "aaa", "bbb", "ccc", "ddd"
	grid1 := GetNewGrid("")
	grid1.OwnerUuid = &uuid1
	grid1.DefaultAccessUuid = &uuid2
	grid1.ViewAccessUuid = &uuid3
	grid1.EditAccessUuid = &uuid4
	grid2 := GetNewGrid("")
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
	grid := GetNewGrid("")
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

func TestGetRowsColumnDefinitions(t *testing.T) {
	grid := GetNewGrid("xxx")
	got := grid.GetRowsColumnDefinitions()
	expect := ", text1 text, text2 text, text3 text, text4 text, text5 text, text6 text, text7 text, text8 text, text9 text, text10 text, int1 integer, int2 integer, int3 integer, int4 integer, int5 integer, int6 integer, int7 integer, int8 integer, int9 integer, int10 integer"
	if got != expect {
		t.Errorf(`Got %s instead of %s.`, got, expect)
	}
}

func TestGetRowsColumnDefinitionsForExportDb(t *testing.T) {
	tests := []struct {
		test   string
		uuid   string
		expect string
	}{
		{"1", UuidGrids, ", text1, text2, text3"},
		{"2", UuidColumns, ", text1, text2, text3, int1"},
		{"3", UuidRelationships, ", text1, text2, text3, text4, text5"},
		{"4", UuidMigrations, ", text1, int1"},
		{"5", UuidUsers, ", text1, text2, text3, text4"},
		{"6", UuidTransactions, ", text1"},
		{"7", "xxx", ", text1, text2, text3, text4, text5, text6, text7, text8, text9, text10, int1, int2, int3, int4, int5, int6, int7, int8, int9, int10"},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid(tt.uuid)
			got := grid.getRowsColumnDefinitionsForExportDb()
			if got != tt.expect {
				t.Errorf(`Got %s instead of %s.`, got, tt.expect)
			}
		})
	}
}

func TestGetRowsQueryForExportDb(t *testing.T) {
	grid := GetNewGrid("xxx")
	got := grid.GetRowsQueryForExportDb()
	expect := "SELECT uuid, gridUuid, created, createdBy, updated, updatedBy, text1, text2, text3, text4, text5, text6, text7, text8, text9, text10, int1, int2, int3, int4, int5, int6, int7, int8, int9, int10, enabled, revision FROM rows ORDER BY created"
	if got != expect {
		t.Errorf(`Got %s instead of %s.`, got, expect)
	}
}

func TestGetRowsQueryForSeedData(t *testing.T) {
	grid := GetNewGrid("xxx")
	got := grid.GetRowsQueryForSeedData()
	expect := "SELECT uuid, revision FROM rows WHERE gridUuid = $1 AND uuid = $2"
	if got != expect {
		t.Errorf(`Got %s instead of %s.`, got, expect)
	}
}

func TestGetInsertStatementForSeedRowDb(t *testing.T) {
	grid := GetNewGrid("xxx")
	got := grid.GetInsertStatementForSeedRowDb()
	expect := "INSERT INTO rows (uuid, revision, created, updated, createdBy, updatedBy, enabled, gridUuid, text1, text2, text3, text4, text5, text6, text7, text8, text9, text10, int1, int2, int3, int4, int5, int6, int7, int8, int9, int10) VALUES ($1, 1, NOW(), NOW(), $2, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)"
	if got != expect {
		t.Errorf(`Got %s instead of %s.`, got, expect)
	}
}

func TestGetInsertStatementParametersForSeedRowDb(t *testing.T) {
	tests := []struct {
		test   string
		uuid   string
		expect string
	}{
		{"1", UuidGrids, ", $5, $6, $7"},
		{"2", UuidColumns, ", $5, $6, $7, $8"},
		{"3", UuidRelationships, ", $5, $6, $7, $8, $9"},
		{"4", UuidMigrations, ", $5, $6"},
		{"5", UuidUsers, ", $5, $6, $7, $8"},
		{"6", UuidTransactions, ", $5"},
		{"7", "xxx", ", $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24"},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid(tt.uuid)
			got := grid.getInsertStatementParametersForSeedRowDb()
			if got != tt.expect {
				t.Errorf(`Got %s instead of %s.`, got, tt.expect)
			}
		})
	}
}

func TestGetInsertValuesForSeedRowDb(t *testing.T) {
	grid := GetNewGrid(UuidGrids)
	text1 := "yyy"
	row := Row{
		Uuid:     "zzz",
		GridUuid: UuidGrids,
		Enabled:  true,
		Text1:    &text1,
		Text2:    &text1,
		Text3:    &text1,
	}
	got := grid.GetInsertValuesForSeedRowDb("xxx", &row)
	expect := []any{
		"zzz",
		"xxx",
		true,
		UuidGrids,
		&text1,
		&text1,
		&text1,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesMigrationsForSeedRowDb(t *testing.T) {
	grid := GetNewGrid(UuidMigrations)
	text1 := "yyy"
	int1 := int64(1)
	row := Row{
		Uuid:     "zzz",
		GridUuid: UuidMigrations,
		Enabled:  true,
		Text1:    &text1,
		Int1:     &int1,
	}
	got := grid.GetInsertValuesForSeedRowDb("xxx", &row)
	expect := []any{
		"zzz",
		"xxx",
		true,
		UuidMigrations,
		&text1,
		&int1,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesUsersForSeedRowDb(t *testing.T) {
	grid := GetNewGrid(UuidUsers)
	text1 := "yyy"
	text2 := "aaa"
	text3 := "bbb"
	text4 := "ccc"
	row := Row{
		Uuid:     "zzz",
		GridUuid: UuidUsers,
		Enabled:  true,
		Text1:    &text1,
		Text2:    &text2,
		Text3:    &text3,
		Text4:    &text4,
	}
	got := grid.GetInsertValuesForSeedRowDb("xxx", &row)
	expect := []any{
		"zzz",
		"xxx",
		true,
		UuidUsers,
		&text1,
		&text2,
		&text3,
		&text4,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetInsertValuesTrasnactionsForSeedRowDb(t *testing.T) {
	grid := GetNewGrid(UuidTransactions)
	text1 := "yyy"
	row := Row{
		Uuid:     "zzz",
		GridUuid: UuidTransactions,
		Enabled:  true,
		Text1:    &text1,
	}
	got := grid.GetInsertValuesForSeedRowDb("xxx", &row)
	expect := []any{
		"zzz",
		"xxx",
		true,
		UuidTransactions,
		&text1,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateValuesForSeedRowDb(t *testing.T) {
	grid := GetNewGrid(UuidGrids)
	text1 := "yyy"
	const revision int8 = 10
	row := Row{
		Uuid:     "zzz",
		GridUuid: UuidGrids,
		Text1:    &text1,
		Text2:    &text1,
		Text3:    &text1,
		Revision: revision,
		Enabled:  true,
	}
	got := grid.GetUpdateValuesForSeedRowDb("xxx", &row)
	expect := []any{
		UuidGrids,
		"zzz",
		revision,
		"xxx",
		true,
		&text1,
		&text1,
		&text1,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetUpdateStatementParametersForSeedRowDb(t *testing.T) {
	tests := []struct {
		test   string
		uuid   string
		expect string
	}{
		{"1", UuidGrids, ", text1 = $6, text2 = $7, text3 = $8"},
		{"2", UuidColumns, ", text1 = $6, text2 = $7, text3 = $8, int1 = $9"},
		{"3", UuidRelationships, ", text1 = $6, text2 = $7, text3 = $8, text4 = $9, text5 = $10"},
		{"4", "xxx", ", text1 = $6, text2 = $7, text3 = $8, text4 = $9, text5 = $10, text6 = $11, text7 = $12, text8 = $13, text9 = $14, text10 = $15, int1 = $16, int2 = $17, int3 = $18, int4 = $19, int5 = $20, int6 = $21, int7 = $22, int8 = $23, int9 = $24, int10 = $25"},
	}
	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			grid := GetNewGrid(tt.uuid)
			got := grid.getUpdateStatementParametersForSeedRowDb()
			if got != tt.expect {
				t.Errorf(`Got %s instead of %s.`, got, tt.expect)
			}
		})
	}
}

func TestGetUpdateStatementForSeedRowDb(t *testing.T) {
	grid := GetNewGrid("xxx")
	got := grid.GetUpdateStatementForSeedRowDb()
	expect := "UPDATE rows SET revision = $3, updated = NOW(), updatedBy = $4, enabled = $5, text1 = $6, text2 = $7, text3 = $8, text4 = $9, text5 = $10, text6 = $11, text7 = $12, text8 = $13, text9 = $14, text10 = $15, int1 = $16, int2 = $17, int3 = $18, int4 = $19, int5 = $20, int6 = $21, int7 = $22, int8 = $23, int9 = $24, int10 = $25 WHERE gridUuid = $1 AND uuid = $2"
	if got != expect {
		t.Errorf(`Got %s instead of %s.`, got, expect)
	}
}
