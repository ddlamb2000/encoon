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
	ctxChan     chan ApiResponse
	transaction *model.Row
}

type ApiParameters struct {
	DbName               string
	UserUuid             string
	UserName             string
	GridUuid             string
	Uuid                 string
	filterColumnOwned    bool
	filterColumnName     string
	filterColumnGridUuid string
	filterColumnValue    string
}

type ApiResponse struct {
	Grid                   *model.Grid         `json:"grid"`
	CountRows              int                 `json:"countRows"`
	Rows                   []model.Row         `json:"rows"`
	RowsAdded              []*model.Row        `json:"rowsAdded,omitempty"`
	RowsEdited             []*model.Row        `json:"rowsEdited,omitempty"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted,omitempty"`
	ReferenceValuesAdded   []gridReferencePost `json:"referencedValuesAdded,omitempty"`
	ReferenceValuesRemoved []gridReferencePost `json:"referencedValuesRemoved,omitempty"`
	Err                    error               `json:"err,omitempty"`
	TimeOut                bool                `json:"timeOut,omitempty"`
	System                 bool                `json:"system,omitempty"`
	Forbidden              bool                `json:"forbidden,omitempty"`
	CanViewRows            bool                `json:"canViewRows"`
	CanEditRows            bool                `json:"canEditRows"`
	CanAddRows             bool                `json:"canAddRows"`
	CanEditGrid            bool                `json:"canEditGrid"`
}

type gridPost struct {
	RowsAdded              []*model.Row        `json:"rowsAdded"`
	RowsEdited             []*model.Row        `json:"rowsEdited"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted"`
	ReferenceValuesAdded   []gridReferencePost `json:"referencedValuesAdded"`
	ReferenceValuesRemoved []gridReferencePost `json:"referencedValuesRemoved"`
}

type gridReferencePost struct {
	ColumnName string `json:"columnName"`
	ColumnUuid string `json:"columnUuid"`
	FromUuid   string `json:"fromUuid"`
	ToGridUuid string `json:"toGridUuid"`
	ToUuid     string `json:"uuid"`
	Owned      bool   `json:"owned"`
}
