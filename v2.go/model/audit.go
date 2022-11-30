// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package model

import "time"

type Audit struct {
	ColumnName    string     `json:"columnName" yaml:"columnName"`
	Uuid          string     `json:"uuid" yaml:"uuid"`
	Created       *time.Time `json:"created" yaml:"created"`
	CreatedBy     *string    `json:"createdBy" yaml:"createdBy"`
	CreatedByName *string    `json:"createdByName,omitempty" yaml:"createdByName,omitempty"`
	ActionName    string     `json:"actionName" yaml:"actionName"`
}

func GetNewAudit() *Audit {
	return new(Audit)
}

func (audit *Audit) SetActionName() {
	switch audit.ColumnName {
	case "relationship1":
		audit.ActionName = "Created"
	case "relationship2":
		audit.ActionName = "Updated"
	case "relationship3":
		audit.ActionName = "Deleted"
	}
}
