// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
	CanViewRows bool      `json:"canViewRows"`
	CanEditRows bool      `json:"canEditRows"`
	CanAddRows  bool      `json:"canAddRows"`
	Columns     []*Column `json:"columns,omitempty"`

	OwnerUuid         *string         `json:"-"`
	DefaultAccessUuid *string         `json:"-"`
	ViewAccessUuid    *string         `json:"-"`
	EditAccessUuid    *string         `json:"-"`
	Owners            map[string]bool `json:"-"`
	DefaultAccess     map[string]bool `json:"-"`
	ViewAccess        map[string]bool `json:"-"`
	EditAccess        map[string]bool `json:"-"`
}

func GetNewGrid() *Grid {
	grid := new(Grid)
	grid.Owners = make(map[string]bool)
	grid.DefaultAccess = make(map[string]bool)
	grid.ViewAccess = make(map[string]bool)
	grid.EditAccess = make(map[string]bool)
	return grid
}

func (grid *Grid) SetPathAndDisplayString(dbName string) {
	grid.Path = fmt.Sprintf("/%s/%s", dbName, grid.Uuid)
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

func (grid *Grid) SetViewEditAccessFlags(userUuid string) {
	switch {
	case grid.Owners[userUuid]:
		grid.CanViewRows = true
		grid.CanEditRows = true
		grid.CanAddRows = true
	case grid.EditAccess[userUuid]:
		grid.CanViewRows = true
		grid.CanEditRows = true
		grid.CanAddRows = true
	case grid.ViewAccess[userUuid]:
		grid.CanViewRows = true
		grid.CanEditRows = false
		grid.CanAddRows = false
	case grid.DefaultAccess[UuidAccessLevelWriteAccess]:
		grid.CanViewRows = true
		grid.CanEditRows = true
		grid.CanAddRows = true
	case grid.DefaultAccess[UuidAccessLevelReadAccess]:
		grid.CanViewRows = true
		grid.CanEditRows = false
		grid.CanAddRows = false
	case grid.Uuid == UuidGrids || grid.Uuid == UuidColumns:
		grid.CanViewRows = true
		grid.CanEditRows = true
		grid.CanAddRows = true
	case grid.Uuid == UuidAccessLevel || grid.Uuid == UuidUsers || grid.Uuid == UuidColumnTypes:
		grid.CanViewRows = true
		grid.CanEditRows = false
		grid.CanAddRows = false
	default:
		grid.CanViewRows = false
		grid.CanEditRows = false
	}
}
