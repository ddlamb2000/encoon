// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

type entity struct {
	Uuid    string `json:"uuid"`
	Uri     string `json:"uri"`
	Version int8   `json:"version"`
	Enabled bool   `json:"enabled"`
}

type user struct {
	Id        entity `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}