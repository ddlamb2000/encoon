// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import "fmt"

type Column struct {
	Uuid           string  `json:"uuid,omitempty"`
	Owned          bool    `json:"owned"`
	Label          string  `json:"label,omitempty"`
	Name           string  `json:"name,omitempty"`
	Type           string  `json:"type,omitempty"`
	TypeUuid       string  `json:"typeUuid,omitempty"`
	GridUuid       string  `json:"gridUuid,omitempty"`
	Grid           *Grid   `json:"grid,omitempty"`
	GridPromptUuid *string `json:"gridPromptUuid,omitempty"`
	Bidirectional  *bool   `json:"bidirectional,omitempty"`
}

func GetNewColumn() *Column {
	column := new(Column)
	return column
}

func (column Column) String() string {
	return fmt.Sprintf("%s %q", column.Name, column.Label)
}

func (column *Column) IsOwned() bool {
	return column.Owned
}

func (column *Column) IsReference() bool {
	return column.TypeUuid == UuidReferenceColumnType
}

func (column *Column) IsAttribute() bool {
	return column.TypeUuid != UuidReferenceColumnType
}
