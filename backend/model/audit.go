// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import "time"

type Audit struct {
	ColumnName    string     `json:"columnName"`
	Uuid          string     `json:"uuid"`
	Created       *time.Time `json:"created"`
	CreatedBy     *string    `json:"createdBy"`
	CreatedByName *string    `json:"createdByName,omitempty"`
	ActionName    string     `json:"actionName"`
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
