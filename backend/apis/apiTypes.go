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
	ReferenceValuesAdded   []GridReferencePost `json:"referencedValuesAdded,omitempty"`
	ReferenceValuesRemoved []GridReferencePost `json:"referencedValuesRemoved,omitempty"`
	Err                    error               `json:"err,omitempty"`
	TimeOut                bool                `json:"timeOut,omitempty"`
	System                 bool                `json:"system,omitempty"`
	Forbidden              bool                `json:"forbidden,omitempty"`
	CanViewRows            bool                `json:"canViewRows"`
	CanEditRows            bool                `json:"canEditRows"`
	CanAddRows             bool                `json:"canAddRows"`
	CanEditGrid            bool                `json:"canEditGrid"`
}

type GridPost struct {
	RowsAdded              []*model.Row        `json:"rowsAdded"`
	RowsEdited             []*model.Row        `json:"rowsEdited"`
	RowsDeleted            []*model.Row        `json:"rowsDeleted"`
	ReferenceValuesAdded   []GridReferencePost `json:"referencedValuesAdded"`
	ReferenceValuesRemoved []GridReferencePost `json:"referencedValuesRemoved"`
}

type GridReferencePost struct {
	ColumnName string `json:"columnName"`
	FromUuid   string `json:"fromUuid"`
	ToGridUuid string `json:"toGridUuid"`
	ToUuid     string `json:"uuid"`
	Owned      bool   `json:"owned"`
}
