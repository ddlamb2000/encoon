// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

type Reference struct {
	Owned    bool   `json:"owned"`
	Label    string `json:"label,omitempty"`
	Name     string `json:"name,omitempty"`
	GridUuid string `json:"gridUuid,omitempty"`
	Rows     []Row  `json:"rows,omitempty"`
}

func GetNewReference() *Reference {
	reference := new(Reference)
	return reference
}
