// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import "fmt"

type Row struct {
	Uuid      string  `json:"uuid"`
	Version   int8    `json:"version"`
	Enabled   bool    `json:"enabled"`
	Uri       *string `json:"uri"`
	GridUuid  *string `json:"gridUuid"`
	CreateBy  *string `json:"createdBy"`
	UpdatedBy *string `json:"updatedBy"`
	Text01    *string `json:"text01"`
	Text02    *string `json:"text02"`
	Text03    *string `json:"text03"`
	Text04    *string `json:"text04"`
	Text05    *string `json:"text05"`
	Text06    *string `json:"text06"`
	Text07    *string `json:"text07"`
	Text08    *string `json:"text08"`
	Text09    *string `json:"text09"`
	Text10    *string `json:"text10"`

	Path string `json:"path"`
}

func (row Row) String() string {
	return row.Uuid
}

func (row *Row) SetPath(dbName, gridUri string) {
	row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
}
