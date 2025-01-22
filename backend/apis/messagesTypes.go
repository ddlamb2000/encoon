// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package apis

type requestContent struct {
	Action               string   `json:"action"`
	ActionText           string   `json:"actionText,omitempty"`
	GridUuid             string   `json:"gridUuid,omitempty"`
	ColumnUuid           string   `json:"columnUuid,omitempty"`
	Uuid                 string   `json:"uuid,omitempty"`
	Userid               string   `json:"userId,omitempty"`
	Password             string   `json:"password,omitempty"`
	FilterColumnOwned    bool     `json:"filterColumnOwned,omitempty"`
	FilterColumnName     string   `json:"filterColumnName,omitempty"`
	FilterColumnGridUuid string   `json:"filterColumnGridUuid,omitempty"`
	FilterColumnValue    string   `json:"filterColumnValue,omitempty"`
	DataSet              GridPost `json:"dataSet,omitempty"`
}

type responseContent struct {
	Action      string       `json:"action"`
	ActionText  string       `json:"actionText,omitempty"`
	Status      string       `json:"status"`
	GridUuid    string       `json:"gridUuid,omitempty"`
	ColumnUuid  string       `json:"columnUuid,omitempty"`
	Uuid        string       `json:"uuid,omitempty"`
	TextMessage string       `json:"textMessage,omitempty"`
	FirstName   string       `json:"firstName,omitempty"`
	LastName    string       `json:"lastName,omitempty"`
	JWT         string       `json:"jwt,omitempty"`
	DataSet     GridResponse `json:"dataSet,omitempty"`
}
