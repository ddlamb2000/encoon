// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

type Reference struct {
	Label string `json:"label,omitempty" yaml:"label,omitempty"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Rows  []Row  `json:"rows,omitempty" yaml:"rows,omitempty"`
}

func GetNewReference() *Reference {
	reference := new(Reference)
	return reference
}
