// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
	OwnerUuid         *string   `json:"-"`
	DefaultAccessUuid *string   `json:"-"`
	ViewAccessUuid    *string   `json:"-"`
	EditAccessUuid    *string   `json:"-"`
	CanView           bool      `json:"canViewGrid"`
	CanEdit           bool      `json:"canEditGrid"`
	SpecialAccess     bool      `json:"gridSpecialAccess"`
	Columns           []*Column `json:"columns,omitempty"`
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
	case grid.OwnerUuid != nil && *grid.OwnerUuid == userUuid:
		grid.CanView = true
		grid.CanEdit = true
		grid.SpecialAccess = false
	case grid.EditAccessUuid != nil && *grid.EditAccessUuid == userUuid:
		grid.CanView = true
		grid.CanEdit = true
		grid.SpecialAccess = false
	case grid.ViewAccessUuid != nil && *grid.ViewAccessUuid == userUuid:
		grid.CanView = true
		grid.CanEdit = false
		grid.SpecialAccess = false
	case grid.DefaultAccessUuid != nil && *grid.DefaultAccessUuid == UuidAccessLevelWriteAccess:
		grid.CanView = true
		grid.CanEdit = true
		grid.SpecialAccess = false
	case grid.DefaultAccessUuid != nil && *grid.DefaultAccessUuid == UuidAccessLevelReadAccess:
		grid.CanView = true
		grid.CanEdit = false
		grid.SpecialAccess = false
	case grid.DefaultAccessUuid != nil && *grid.DefaultAccessUuid == UuidAccessLevelSpecialAccess:
		grid.CanView = false
		grid.CanEdit = false
		grid.SpecialAccess = true
	default:
		grid.CanView = false
		grid.CanEdit = false
		grid.SpecialAccess = false
	}
}
