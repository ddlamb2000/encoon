// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

type Column struct {
	Label string `json:"label,omitempty" yaml:"label,omitempty"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
}
