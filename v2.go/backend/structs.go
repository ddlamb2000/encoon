// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

type entity struct {
	Uuid    string `json:"uuid"`
	Version int8   `json:"version"`
	Enabled bool   `json:"enabled"`
}
