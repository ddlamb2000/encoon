// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
	Columns []*Column `json:"columns,omitempty" yaml:"columns,omitempty"`
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
	default:
		return "rows"
	}
}
