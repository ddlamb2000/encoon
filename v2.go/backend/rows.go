// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import "fmt"

const numberOfTextFields = 10
const numberOfIntFields = 10

type Row struct {
	Uuid      string  `json:"uuid"`
	Version   int8    `json:"version"`
	Enabled   bool    `json:"enabled"`
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
	Int01     *string `json:"int01"`
	Int02     *string `json:"int02"`
	Int03     *string `json:"int03"`
	Int04     *string `json:"int04"`
	Int05     *string `json:"int05"`
	Int06     *string `json:"int06"`
	Int07     *string `json:"int07"`
	Int08     *string `json:"int08"`
	Int09     *string `json:"int09"`
	Int10     *string `json:"int10"`

	Path string `json:"path"`
}

func (row Row) String() string {
	return row.Uuid
}

func (row *Row) SetPath(dbName, gridUri string) {
	row.Path = fmt.Sprintf("/%s/%s/%s", dbName, gridUri, row.Uuid)
}

func getRowsColumnDefinitions() string {
	var columnDefinitions = ""
	for i := 1; i <= numberOfTextFields; i++ {
		columnDefinitions += fmt.Sprintf("text%02d text, ", i)
	}
	for i := 1; i <= numberOfIntFields; i++ {
		columnDefinitions += fmt.Sprintf("int%02d integer, ", i)
	}
	return columnDefinitions
}
