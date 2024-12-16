// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Column struct {
	Uuid           string  `json:"uuid,omitempty"`
	OrderNumber    *int64  `json:"orderNumber,omitempty"`
	Owned          bool    `json:"owned"`
	Label          *string `json:"label,omitempty"`
	Name           *string `json:"name,omitempty"`
	Type           *string `json:"type,omitempty"`
	TypeUuid       *string `json:"typeUuid,omitempty"`
	GridUuid       *string `json:"gridUuid,omitempty"`
	Grid           *Grid   `json:"grid,omitempty"`
	GridPromptUuid *string `json:"gridPromptUuid,omitempty"`
	Bidirectional  *bool   `json:"bidirectional,omitempty"`
}

func GetNewColumn() *Column {
	column := new(Column)
	return column
}

func (column Column) String() string {
	if column.Name != nil && column.Label != nil {
		return fmt.Sprintf("%s %q", *column.Name, *column.Label)
	}
	return ""
}

func (column *Column) IsOwned() bool {
	return column.Owned
}

func (column *Column) IsReference() bool {
	return column.TypeUuid != nil && *column.TypeUuid == UuidReferenceColumnType
}

func (column *Column) IsAttribute() bool {
	return column.TypeUuid == nil || *column.TypeUuid != UuidReferenceColumnType
}

func (column *Column) GetColumnNamePrefixFromType() string {
	if column.TypeUuid == nil {
		return ""
	}
	switch *column.TypeUuid {
	case UuidIntColumnType:
		return "int"
	case UuidReferenceColumnType:
		return "relationship"
	default:
		return "text"
	}
}

func (column *Column) GetColumnNamePrefixAndIndex() (string, int64) {
	if column.Name != nil {
		prefixes := []string{"int", "relationship", "text"}
		name := *column.Name
		for _, prefix := range prefixes {
			if strings.HasPrefix(name, prefix) {
				indexStr := name[len(prefix):]
				index, err := strconv.Atoi(indexStr)
				if err == nil {
					return prefix, (int64)(index)
				}
			}
		}
	}
	return "", 0
}
