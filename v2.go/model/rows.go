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
	Text01    *string `json:"text01,omitempty" yaml:"text01,omitempty"`
	Text02    *string `json:"text02,omitempty" yaml:"text02,omitempty"`
	Text03    *string `json:"text03,omitempty" yaml:"text03,omitempty"`
	Text04    *string `json:"text04,omitempty" yaml:"text04,omitempty"`
	Text05    *string `json:"text05,omitempty" yaml:"text05,omitempty"`
	Text06    *string `json:"text06,omitempty" yaml:"text06,omitempty"`
	Text07    *string `json:"text07,omitempty" yaml:"text07,omitempty"`
	Text08    *string `json:"text08,omitempty" yaml:"text08,omitempty"`
	Text09    *string `json:"text09,omitempty" yaml:"text09,omitempty"`
	Text10    *string `json:"text10,omitempty" yaml:"text10,omitempty"`
	Int01     *int64  `json:"int01,omitempty" yaml:"int01,omitempty"`
	Int02     *int64  `json:"int02,omitempty" yaml:"int01,omitempty"`
	Int03     *int64  `json:"int03,omitempty" yaml:"int01,omitempty"`
	Int04     *int64  `json:"int04,omitempty" yaml:"int01,omitempty"`
	Int05     *int64  `json:"int05,omitempty" yaml:"int01,omitempty"`
	Int06     *int64  `json:"int06,omitempty" yaml:"int01,omitempty"`
	Int07     *int64  `json:"int07,omitempty" yaml:"int01,omitempty"`
	Int08     *int64  `json:"int08,omitempty" yaml:"int01,omitempty"`
	Int09     *int64  `json:"int09,omitempty" yaml:"int01,omitempty"`
	Int10     *int64  `json:"int10,omitempty" yaml:"int01,omitempty"`
	Path      string  `json:"path" yaml:"path"`
}

func (row Row) String() string {
	return row.Uuid
}

func (row *Row) SetPath(dbName, gridUri string) {
	row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
}
