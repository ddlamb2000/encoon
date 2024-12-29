// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"d.lambert.fr/encoon/apis"
)

type requestContent struct {
	Action     string `json:"action"`
	GridUuid   string `json:"gridUuid,omitempty"`
	ColumnUuid string `json:"columnUuid,omitempty"`
	RowUuid    string `json:"rowUuid,omitempty"`
	Uuid       string `json:"uuid,omitempty"`
	Userid     string `json:"userId,omitempty"`
	Password   string `json:"password,omitempty"`
}

type responseContent struct {
	Status      string           `json:"status"`
	Action      string           `json:"action"`
	GridUuid    string           `json:"gridUuid,omitempty"`
	ColumnUuid  string           `json:"columnUuid,omitempty"`
	RowUuid     string           `json:"rowUuid,omitempty"`
	Uuid        string           `json:"uuid,omitempty"`
	TextMessage string           `json:"textMessage,omitempty"`
	FirstName   string           `json:"firstName,omitempty"`
	LastName    string           `json:"lastName,omitempty"`
	JWT         string           `json:"jwt,omitempty"`
	DataSet     apis.ApiResponse `json:"dataSet,omitempty"`
}

const (
	SuccessStatus = "SUCCESS"
	FailedStatus  = "FAILED"

	ActionAuthentication = "AUTHENTICATION"
	ActionGetGrid        = "LOAD"
	ActionLocateGrid     = "LOCATE"
)
