// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

import (
	"d.lambert.fr/encoon/apis"
)

type requestContent struct {
	Action   string `json:"action"`
	GridUuid string `json:"griduuid,omitempty"`
	Uuid     string `json:"uuid,omitempty"`
	Userid   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
}

type responseContent struct {
	Status      string           `json:"status"`
	Action      string           `json:"action"`
	GridUuid    string           `json:"griduuid,omitempty"`
	Uuid        string           `json:"uuid,omitempty"`
	TextMessage string           `json:"textmessage,omitempty"`
	FirstName   string           `json:"firstname,omitempty"`
	LastName    string           `json:"lastname,omitempty"`
	JWT         string           `json:"jwt,omitempty"`
	DataSet     apis.ApiResponse `json:"dataset,omitempty"`
}

const (
	SuccessStatus = "SUCCESS"
	FailedStatus  = "FAILED"

	ActionAuthentication = "AUTHENTICATION"
	ActionGetGrid        = "LOAD"
)
