// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Row struct {
	Uuid      string  `json:"uuid" yaml:"uuid"`
	Version   int8    `json:"version" yaml:"version"`
	Enabled   bool    `json:"enabled" yaml:"enabled"`
	GridUuid  *string `json:"gridUuid" yaml:"gridUuid"`
	CreateBy  *string `json:"createdBy" yaml:"createdBy"`
	UpdatedBy *string `json:"updatedBy" yaml:"updatedBy"`
	Text1     *string `json:"text1,omitempty" yaml:"text1,omitempty"`
	Text2     *string `json:"text2,omitempty" yaml:"text2,omitempty"`
	Text3     *string `json:"text3,omitempty" yaml:"text3,omitempty"`
	Text4     *string `json:"text4,omitempty" yaml:"text4,omitempty"`
	Text5     *string `json:"text5,omitempty" yaml:"text5,omitempty"`
	Text6     *string `json:"text6,omitempty" yaml:"text6,omitempty"`
	Text7     *string `json:"text7,omitempty" yaml:"text7,omitempty"`
	Text8     *string `json:"text8,omitempty" yaml:"text8,omitempty"`
	Text9     *string `json:"text9,omitempty" yaml:"text9,omitempty"`
	Text10    *string `json:"text10,omitempty" yaml:"text10,omitempty"`
	Int1      *int64  `json:"int1,omitempty" yaml:"int1,omitempty"`
	Int2      *int64  `json:"int2,omitempty" yaml:"int1,omitempty"`
	Int3      *int64  `json:"int3,omitempty" yaml:"int1,omitempty"`
	Int4      *int64  `json:"int4,omitempty" yaml:"int1,omitempty"`
	Int5      *int64  `json:"int5,omitempty" yaml:"int1,omitempty"`
	Int6      *int64  `json:"int6,omitempty" yaml:"int1,omitempty"`
	Int7      *int64  `json:"int7,omitempty" yaml:"int1,omitempty"`
	Int8      *int64  `json:"int8,omitempty" yaml:"int1,omitempty"`
	Int9      *int64  `json:"int9,omitempty" yaml:"int1,omitempty"`
	Int10     *int64  `json:"int10,omitempty" yaml:"int1,omitempty"`
	Path      string  `json:"path" yaml:"path"`
}

func (row Row) String() string {
	return row.Uuid
}

func (row *Row) SetPath(dbName, gridUri string) {
	row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
}
