// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
}

func (row *Grid) SetPath(dbName, gridUri string) {
	row.Path = fmt.Sprintf("/%s/%s", dbName, gridUri)
}
