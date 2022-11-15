// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "fmt"

type Grid struct {
	Row
}

func (grid *Grid) SetPath(dbName, gridUri string) {
	grid.Path = fmt.Sprintf("/%s/%s", dbName, gridUri)
}
