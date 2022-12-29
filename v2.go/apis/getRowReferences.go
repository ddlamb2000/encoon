// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"d.lambert.fr/encoon/model"
)

func getRelationshipsForRow(r apiRequestParameters, grid *model.Grid, row *model.Row) error {
	t := r.startTiming()
	defer r.stopTiming("getRelationshipsForRow()", t)
	r.trace("getRelationshipsForRow(%s, %s)", grid, row)
	for _, col := range grid.Columns {
		var referencedRows []model.Row
		var err error
		if col.IsReference() {
			r.trace("getRelationshipsForRow() - col=%s", col)
			r.log("getRelationshipsForRow() - col.Name=%s, col.GridUuid=%s", col.Name, col.GridUuid)
			referencedRows, err = getReferencedRowsForRow(r, row, col.Name, col.GridUuid, col.IsOwned())
			if err != nil {
				return r.logAndReturnError("Error when retrieving referenced rows: %v.", err)
			}
		}
		if len(referencedRows) > 0 {
			reference := model.GetNewReference()
			reference.Owned = col.Owned
			reference.Name = col.Name
			reference.Label = col.Label
			reference.GridUuid = col.GridUuid
			reference.Rows = referencedRows
			row.References = append(row.References, reference)
		}
	}
	return nil
}

func getReferencedRowsForRow(r apiRequestParameters, parentRow *model.Row, referenceName, referenceColumnUuid string, owned bool) ([]model.Row, error) {
	t := r.startTiming()
	defer r.stopTiming("getReferencedRowsForRow()", t)
	query := getQueryReferencedRowsForRow(owned)
	parms := getQueryParametersReferencedRowsForRow(referenceName, referenceColumnUuid, parentRow)
	r.trace("getReferencedRowsForRow(%s, %s) - query=%s ; parms=%s", parentRow, referenceName, query, parms)
	rows, err := r.db.QueryContext(r.ctx, query, parms...)
	if err != nil {
		return nil, r.logAndReturnError("Error when querying referenced rows: %v.", err)
	}
	defer rows.Close()
	rowSet := make([]model.Row, 0)
	for rows.Next() {
		var referencedUuid string
		var referencedGridUuid string
		if err := rows.Scan(&referencedGridUuid, &referencedUuid); err != nil {
			return nil, r.logAndReturnError("Error when scanning referenced rows: %v.", err)
		}
		grid, _ := getGridForGridsApi(r, referencedGridUuid)
		if grid != nil {
			rows, _, err := getRowSetForGridsApi(r, grid, referencedUuid, false, false)
			if err == nil {
				rowSet = append(rowSet, rows...)
			}
		}
	}
	return rowSet, nil
}

// function is available for mocking
var getQueryReferencedRowsForRow = func(owned bool) string {
	if owned {
		return "SELECT text4, " +
			"text5 " +
			"FROM relationships " +
			"WHERE gridUuid = $1 " +
			"AND text1 = $2 " +
			"AND text2 = $3 " +
			"AND text3 = $4 " +
			"AND text2 = $5 " +
			"AND enabled = true"
	} else {
		return "SELECT text2, " +
			"text3 " +
			"FROM relationships " +
			"WHERE gridUuid = $1 " +
			"AND text1 = $2 " +
			"AND text4 = $3 " +
			"AND text5 = $4 " +
			"AND text2 = $5 " +
			"AND enabled = true"
	}
}

func getQueryParametersReferencedRowsForRow(referenceName, referenceColumnUuid string, parentRow *model.Row) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, referenceName)
	parameters = append(parameters, parentRow.GridUuid)
	parameters = append(parameters, parentRow.Uuid)
	parameters = append(parameters, referenceColumnUuid)
	return parameters
}
