// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2024

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/model"
)

type ApiRequest struct {
	ctx         context.Context
	p           ApiParameters
	db          *sql.DB
	ctxChan     chan GridResponse
	transaction *model.Row
}

type ApiParameters struct {
	DbName               string
	UserUuid             string
	UserName             string
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
	Action               string       `json:"action"`
	ActionText           string       `json:"actionText,omitempty"`
	Status               string       `json:"status"`
	GridUuid             string       `json:"gridUuid,omitempty"`
	ColumnUuid           string       `json:"columnUuid,omitempty"`
	Uuid                 string       `json:"uuid,omitempty"`
	TextMessage          string       `json:"textMessage,omitempty"`
	FirstName            string       `json:"firstName,omitempty"`
	LastName             string       `json:"lastName,omitempty"`
	JWT                  string       `json:"jwt,omitempty"`
	FilterColumnOwned    bool         `json:"filterColumnOwned,omitempty"`
	FilterColumnName     string       `json:"filterColumnName,omitempty"`
	FilterColumnGridUuid string       `json:"filterColumnGridUuid,omitempty"`
	FilterColumnValue    string       `json:"filterColumnValue,omitempty"`
	DataSet              GridResponse `json:"dataSet,omitempty"`
}

type GridResponse struct {
	Grid                   *model.Grid         `json:"grid,omitempty"`
	CountRows              int                 `json:"countRows,omitempty"`
	Rows                   []model.Row         `json:"rows,omitempty"`
	RowsAdded              []*model.Row        `json:"rowsAdded,omitempty"`
	RowsEdited             []*model.Row        `json:"rowsEdited,omitempty"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted,omitempty"`
	ReferenceValuesAdded   []GridReferencePost `json:"referencedValuesAdded,omitempty"`
	ReferenceValuesRemoved []GridReferencePost `json:"referencedValuesRemoved,omitempty"`
	Err                    error               `json:"err,omitempty"`
	TimeOut                bool                `json:"timeOut,omitempty"`
	System                 bool                `json:"system,omitempty"`
	Forbidden              bool                `json:"forbidden,omitempty"`
	CanViewRows            bool                `json:"canViewRows,omitempty"`
	CanEditRows            bool                `json:"canEditRows,omitempty"`
	CanAddRows             bool                `json:"canAddRows,omitempty"`
	CanEditGrid            bool                `json:"canEditGrid,omitempty"`
}

type GridPost struct {
	RowsAdded              []*model.Row        `json:"rowsAdded,omitempty"`
	RowsEdited             []*model.Row        `json:"rowsEdited,omitempty"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted,omitempty"`
	ReferenceValuesAdded   []GridReferencePost `json:"referencedValuesAdded,omitempty"`
	ReferenceValuesRemoved []GridReferencePost `json:"referencedValuesRemoved,omitempty"`
}

type GridReferencePost struct {
	ColumnName string `json:"columnName"`
	FromUuid   string `json:"fromUuid"`
	ToGridUuid string `json:"toGridUuid"`
	ToUuid     string `json:"uuid"`
	Owned      bool   `json:"owned"`
}
