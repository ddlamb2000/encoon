// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import (
	"fmt"
	"time"
)

type Row struct {
	GridUuid      string       `json:"gridUuid" yaml:"gridUuid"`
	Uuid          string       `json:"uuid" yaml:"uuid"`
	Revision      int8         `json:"revision" yaml:"revision"`
	Enabled       bool         `json:"enabled" yaml:"enabled"`
	Created       *time.Time   `json:"created" yaml:"created"`
	CreatedBy     *string      `json:"createdBy" yaml:"createdBy"`
	Updated       *time.Time   `json:"updated" yaml:"updated"`
	UpdatedBy     *string      `json:"updatedBy" yaml:"updatedBy"`
	Text1         *string      `json:"text1,omitempty" yaml:"text1,omitempty"`
	Text2         *string      `json:"text2,omitempty" yaml:"text2,omitempty"`
	Text3         *string      `json:"text3,omitempty" yaml:"text3,omitempty"`
	Text4         *string      `json:"text4,omitempty" yaml:"text4,omitempty"`
	Text5         *string      `json:"text5,omitempty" yaml:"text5,omitempty"`
	Text6         *string      `json:"text6,omitempty" yaml:"text6,omitempty"`
	Text7         *string      `json:"text7,omitempty" yaml:"text7,omitempty"`
	Text8         *string      `json:"text8,omitempty" yaml:"text8,omitempty"`
	Text9         *string      `json:"text9,omitempty" yaml:"text9,omitempty"`
	Text10        *string      `json:"text10,omitempty" yaml:"text10,omitempty"`
	Int1          *int64       `json:"int1,omitempty" yaml:"int1,omitempty"`
	Int2          *int64       `json:"int2,omitempty" yaml:"int2,omitempty"`
	Int3          *int64       `json:"int3,omitempty" yaml:"int3,omitempty"`
	Int4          *int64       `json:"int4,omitempty" yaml:"int4,omitempty"`
	Int5          *int64       `json:"int5,omitempty" yaml:"int5,omitempty"`
	Int6          *int64       `json:"int6,omitempty" yaml:"int6,omitempty"`
	Int7          *int64       `json:"int7,omitempty" yaml:"int7,omitempty"`
	Int8          *int64       `json:"int8,omitempty" yaml:"int8,omitempty"`
	Int9          *int64       `json:"int9,omitempty" yaml:"int9,omitempty"`
	Int10         *int64       `json:"int10,omitempty" yaml:"int10,omitempty"`
	Path          string       `json:"path,omitempty" yaml:"path,omitempty"`
	DisplayString string       `json:"displayString,omitempty" yaml:"displayString,omitempty"`
	References    []*Reference `json:"references,omitempty" yaml:"references,omitempty"`
	TmpUuid       string
}

func (row Row) String() string {
	return row.Uuid
}

func (row *Row) SetPathAndDisplayString(dbName string) {
	row.Path = fmt.Sprintf("/%s/%s/%s", dbName, row.GridUuid, row.Uuid)
	if row.Text1 != nil {
		row.DisplayString = fmt.Sprintf("%s", *row.Text1)
	}
}
