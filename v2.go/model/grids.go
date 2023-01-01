// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
	Columns []*Column `json:"columns,omitempty"`
	Usages  []*Column `json:"columnsUsage,omitempty"`

	Owners            map[string]bool `json:"-"`
	DefaultAccess     map[string]bool `json:"-"`
	ViewAccess        map[string]bool `json:"-"`
	EditAccess        map[string]bool `json:"-"`
	OwnerUuid         *string         `json:"-"`
	DefaultAccessUuid *string         `json:"-"`
	ViewAccessUuid    *string         `json:"-"`
	EditAccessUuid    *string         `json:"-"`
}

func GetNewGrid() *Grid {
	grid := new(Grid)
	grid.Owners = make(map[string]bool)
	grid.DefaultAccess = make(map[string]bool)
	grid.ViewAccess = make(map[string]bool)
	grid.EditAccess = make(map[string]bool)
	return grid
}

func (grid *Grid) SetDisplayString(dbName string) {
	if grid.Text1 != nil {
		grid.DisplayString = fmt.Sprintf("%s", *grid.Text1)
	}
}

func (grid *Grid) GetTableName() string {
	switch grid.Uuid {
	case UuidGrids:
		return "grids"
	case UuidColumns:
		return "columns"
	case UuidRelationships:
		return "relationships"
	case UuidMigrations:
		return "migrations"
	case UuidUsers:
		return "users"
	case UuidTransactions:
		return "transactions"
	default:
		return "rows"
	}
}

func (grid *Grid) CopyAccessToOtherGrid(otherGrid *Grid) {
	if grid.OwnerUuid != nil {
		otherGrid.Owners[*grid.OwnerUuid] = true
	}
	if grid.DefaultAccessUuid != nil {
		otherGrid.DefaultAccess[*grid.DefaultAccessUuid] = true
	}
	if grid.ViewAccessUuid != nil {
		otherGrid.ViewAccess[*grid.ViewAccessUuid] = true
	}
	if grid.EditAccessUuid != nil {
		otherGrid.EditAccess[*grid.EditAccessUuid] = true
	}
}

func (grid *Grid) GetViewEditAccessFlags(userUuid string) (canViewRows, canEditRows, canEditOwnedRows, canAddRows bool) {
	switch {
	case grid.Owners[userUuid] && grid.Uuid == UuidTransactions:
		return true, false, false, false
	case grid.Owners[userUuid]:
		return true, true, true, true
	case grid.EditAccess[userUuid] || grid.DefaultAccess[UuidAccessLevelWriteAccess]:
		return true, false, true, true
	case grid.ViewAccess[userUuid] || grid.DefaultAccess[UuidAccessLevelReadAccess]:
		return true, false, false, false
	case grid.Uuid == UuidGrids || grid.Uuid == UuidRelationships || grid.Uuid == UuidColumns:
		return true, false, true, true
	case grid.Uuid == UuidAccessLevels || grid.Uuid == UuidUsers || grid.Uuid == UuidColumnTypes:
		return true, false, false, false
	}
	return false, false, false, false
}

func (grid *Grid) HasOwnership(userUuid string) bool {
	return grid.Owners[userUuid]
}
