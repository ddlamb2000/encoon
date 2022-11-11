// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import "fmt"

type Row struct {
	Uuid      string  `json:"uuid"`
	Version   int8    `json:"version"`
	Enabled   bool    `json:"enabled"`
	GridUuid  *string `json:"gridUuid"`
	CreateBy  *string `json:"createdBy"`
	UpdatedBy *string `json:"updatedBy"`
	Text01    *string `json:"text01,omitempty"`
	Text02    *string `json:"text02,omitempty"`
	Text03    *string `json:"text03,omitempty"`
	Text04    *string `json:"text04,omitempty"`
	Text05    *string `json:"text05,omitempty"`
	Text06    *string `json:"text06,omitempty"`
	Text07    *string `json:"text07,omitempty"`
	Text08    *string `json:"text08,omitempty"`
	Text09    *string `json:"text09,omitempty"`
	Text10    *string `json:"text10,omitempty"`
	Int01     *int64  `json:"int01,omitempty"`
	Int02     *int64  `json:"int02,omitempty"`
	Int03     *int64  `json:"int03,omitempty"`
	Int04     *int64  `json:"int04,omitempty"`
	Int05     *int64  `json:"int05,omitempty"`
	Int06     *int64  `json:"int06,omitempty"`
	Int07     *int64  `json:"int07,omitempty"`
	Int08     *int64  `json:"int08,omitempty"`
	Int09     *int64  `json:"int09,omitempty"`
	Int10     *int64  `json:"int10,omitempty"`

	Path string `json:"path"`
}

func (row Row) String() string {
	return row.Uuid
}

func (row *Row) SetPath(dbName, gridUri string) {
	row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
}
