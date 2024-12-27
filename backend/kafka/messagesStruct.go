// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package kafka

type requestContent struct {
	Action   string `json:"action"`
	GridUuid string `json:"griduuid,omitempty"`
	RowUuid  string `json:"rowuuid,omitempty"`
	Userid   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
}

type responseContent struct {
	Status      string `json:"status"`
	Action      string `json:"action"`
	GridUuid    string `json:"griduuid,omitempty"`
	RowUuid     string `json:"rowuuid,omitempty"`
	TextMessage string `json:"textmessage,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	JWT         string `json:"jwt,omitempty"`
}

const (
	SuccessStatus = "SUCCESS"
	FailedStatus  = "FAILED"

	ActionAuthentication = "AUTHENTICATION"
)
