// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
	CanView bool      `json:"canViewGrid"`
	CanEdit bool      `json:"canEditGrid"`
	Columns []*Column `json:"columns,omitempty"`

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

func (grid *Grid) SetViewEditAccessFlags(userUuid string) {
	switch {
	case grid.Owners[userUuid]:
		grid.CanView = true
		grid.CanEdit = true
	case grid.EditAccess[userUuid]:
		grid.CanView = true
		grid.CanEdit = true
	case grid.ViewAccess[userUuid]:
		grid.CanView = true
		grid.CanEdit = false
	case grid.DefaultAccess[UuidAccessLevelWriteAccess]:
		grid.CanView = true
		grid.CanEdit = true
	case grid.DefaultAccess[UuidAccessLevelReadAccess]:
		grid.CanView = true
		grid.CanEdit = false
	default:
		grid.CanView = false
		grid.CanEdit = false
	}
}

func (grid *Grid) GetCanView() bool {
	return grid.CanView || grid.isSpecial()
}

func (grid *Grid) GetCanEdit() bool {
	return grid.CanEdit || grid.isSpecial()
}

func (grid *Grid) isSpecial() bool {
	return grid.Uuid == UuidGrids || grid.Uuid == UuidColumns
}
