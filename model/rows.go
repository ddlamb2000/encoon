// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"fmt"
	"time"

	"d.lambert.fr/encoon/utils"
)

const (
	NumberOfTextFields = 10
	NumberOfIntFields  = 10
)

type Row struct {
	GridUuid      string       `json:"gridUuid"`
	Uuid          string       `json:"uuid"`
	Revision      int8         `json:"revision"`
	Enabled       bool         `json:"enabled"`
	Created       *time.Time   `json:"created"`
	CreatedBy     *string      `json:"createdBy"`
	Updated       *time.Time   `json:"updated"`
	UpdatedBy     *string      `json:"updatedBy"`
	Text1         *string      `json:"text1,omitempty"`
	Text2         *string      `json:"text2,omitempty"`
	Text3         *string      `json:"text3,omitempty"`
	Text4         *string      `json:"text4,omitempty"`
	Text5         *string      `json:"text5,omitempty"`
	Text6         *string      `json:"text6,omitempty"`
	Text7         *string      `json:"text7,omitempty"`
	Text8         *string      `json:"text8,omitempty"`
	Text9         *string      `json:"text9,omitempty"`
	Text10        *string      `json:"text10,omitempty"`
	Int1          *int64       `json:"int1,omitempty"`
	Int2          *int64       `json:"int2,omitempty"`
	Int3          *int64       `json:"int3,omitempty"`
	Int4          *int64       `json:"int4,omitempty"`
	Int5          *int64       `json:"int5,omitempty"`
	Int6          *int64       `json:"int6,omitempty"`
	Int7          *int64       `json:"int7,omitempty"`
	Int8          *int64       `json:"int8,omitempty"`
	Int9          *int64       `json:"int9,omitempty"`
	Int10         *int64       `json:"int10,omitempty"`
	DisplayString string       `json:"displayString,omitempty"`
	CanViewRow    bool         `json:"canViewRow"`
	CanEditRow    bool         `json:"canEditRow"`
	References    []*Reference `json:"references,omitempty"`
	Audits        []*Audit     `json:"audits,omitempty"`

	TmpUuid string `json:"TmpUuid,omitempty"`
}

func GetNewRow() *Row {
	return new(Row)
}

func GetNewRowWithUuid() *Row {
	row := GetNewRow()
	row.Uuid = utils.GetNewUUID()
	return row
}

func (row *Row) SetViewEditAccessFlags(grid *Grid, userUuid string) {
	if grid == nil {
		if *row.CreatedBy == userUuid {
			row.CanViewRow = true
			row.CanEditRow = true
		}
	} else {
		canViewRows, canEditRows, canEditOwnedRows, _ := grid.GetViewEditAccessFlags(userUuid)
		row.CanViewRow = canViewRows
		if canEditRows {
			row.CanEditRow = true
		} else if canEditOwnedRows {
			row.CanEditRow = grid.HasOwnership(userUuid) || *row.CreatedBy == userUuid
		}
	}
}

func (row Row) String() string {
	return row.DisplayString + " [" + row.Uuid + "]"
}

func (row *Row) SetDisplayString(dbName string) {
	if row.Text1 != nil && *row.Text1 != "" {
		row.DisplayString = *row.Text1
	} else if row.Int1 != nil {
		row.DisplayString = fmt.Sprint(*row.Int1)
	} else {
		row.DisplayString = row.Uuid
	}
	if !row.Enabled {
		row.DisplayString += " [DELETED]"
	}
}

func (row *Row) GetRowsQueryOutput() []any {
	output := make([]any, 0)
	output = append(output, &row.Uuid)
	output = append(output, &row.GridUuid)
	output = append(output, &row.Created)
	output = append(output, &row.CreatedBy)
	output = append(output, &row.Updated)
	output = append(output, &row.UpdatedBy)
	switch row.GridUuid {
	case UuidGrids:
		output = append(output, &row.Text1)
		output = append(output, &row.Text2)
		output = append(output, &row.Text3)
	case UuidColumns:
		output = append(output, &row.Text1)
		output = append(output, &row.Text2)
		output = append(output, &row.Text3)
		output = append(output, &row.Int1)
	case UuidRelationships:
		output = append(output, &row.Text1)
		output = append(output, &row.Text2)
		output = append(output, &row.Text3)
		output = append(output, &row.Text4)
		output = append(output, &row.Text5)
	default:
		output = append(output, &row.Text1)
		output = append(output, &row.Text2)
		output = append(output, &row.Text3)
		output = append(output, &row.Text4)
		output = append(output, &row.Text5)
		output = append(output, &row.Text6)
		output = append(output, &row.Text7)
		output = append(output, &row.Text8)
		output = append(output, &row.Text9)
		output = append(output, &row.Text10)
		output = append(output, &row.Int1)
		output = append(output, &row.Int2)
		output = append(output, &row.Int3)
		output = append(output, &row.Int4)
		output = append(output, &row.Int5)
		output = append(output, &row.Int6)
		output = append(output, &row.Int7)
		output = append(output, &row.Int8)
		output = append(output, &row.Int9)
		output = append(output, &row.Int10)
	}
	output = append(output, &row.Enabled)
	output = append(output, &row.Revision)
	return output
}

func (row *Row) AppendRowValuesForSeedRowDb(values []any) []any {
	switch row.GridUuid {
	case UuidGrids:
		values = append(values, row.Text1)
		values = append(values, row.Text2)
		values = append(values, row.Text3)
	case UuidColumns:
		values = append(values, row.Text1)
		values = append(values, row.Text2)
		values = append(values, row.Text3)
		values = append(values, row.Int1)
	case UuidRelationships:
		values = append(values, row.Text1)
		values = append(values, row.Text2)
		values = append(values, row.Text3)
		values = append(values, row.Text4)
		values = append(values, row.Text5)
	default:
		values = append(values, row.Text1)
		values = append(values, row.Text2)
		values = append(values, row.Text3)
		values = append(values, row.Text4)
		values = append(values, row.Text5)
		values = append(values, row.Text6)
		values = append(values, row.Text7)
		values = append(values, row.Text8)
		values = append(values, row.Text9)
		values = append(values, row.Text10)
		values = append(values, row.Int1)
		values = append(values, row.Int2)
		values = append(values, row.Int3)
		values = append(values, row.Int4)
		values = append(values, row.Int5)
		values = append(values, row.Int6)
		values = append(values, row.Int7)
		values = append(values, row.Int8)
		values = append(values, row.Int9)
		values = append(values, row.Int10)
	}
	return values
}
