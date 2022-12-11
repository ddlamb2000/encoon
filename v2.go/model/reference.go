// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

type Reference struct {
	Owned    bool   `json:"owned" yaml:"owned"`
	Label    string `json:"label,omitempty" yaml:"label,omitempty"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	GridUuid string `json:"gridUuid,omitempty" yaml:"gridUuid,omitempty"`
	Rows     []Row  `json:"rows,omitempty" yaml:"rows,omitempty"`
}

func GetNewReference() *Reference {
	reference := new(Reference)
	return reference
}
