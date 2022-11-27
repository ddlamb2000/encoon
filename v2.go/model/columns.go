// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Column struct {
	Label          string  `json:"label,omitempty" yaml:"label,omitempty"`
	Name           string  `json:"name,omitempty" yaml:"name,omitempty"`
	Type           string  `json:"type,omitempty" yaml:"type,omitempty"`
	TypeUuid       string  `json:"typeUuid,omitempty" yaml:"typeUuid,omitempty"`
	GridPromptUuid *string `json:"gridPromptUuid,omitempty" yaml:"gridPromptUuid,omitempty"`
}

func GetNewColumn() *Column {
	column := new(Column)
	return column
}

func (column Column) String() string {
	return fmt.Sprintf("%s %q", column.Name, column.Label)
}

func (column *Column) IsReference() bool {
	return column.TypeUuid == UuidReferenceColumnType
}

func (column *Column) IsAttribute() bool {
	return column.TypeUuid != UuidReferenceColumnType
}
