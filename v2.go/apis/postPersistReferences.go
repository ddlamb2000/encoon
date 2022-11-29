// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type persistGridReferenceDataFunc func(apiRequestParameters, *model.Grid, []*model.Row, gridReferencePost) error

func persistGridReferenceData(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row, refs []gridReferencePost, f persistGridReferenceDataFunc) error {
	for _, ref := range refs {
		err := f(r, grid, addedRows, ref)
		if err != nil {
			_ = r.rollbackTransaction()
			return err
		}
	}
	return nil
}

func postInsertReferenceRow(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost) error {
	query := getInsertStatementForReferenceRow()
	rowUuid := getUuidFromRowsForTmpUuid(r, addedRows, ref.FromUuid)
	parms := getInsertStatementParametersForRefereceRow(r, grid, ref, rowUuid)
	r.trace("postInsertReferenceRow(%s, %v, %v) - query=%s ; parms=%s", grid, addedRows, ref, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Insert referenced row error: %v.", err)
	}
	if grid.Uuid == model.UuidGrids {
		removeGridFromCache(ref.FromUuid)
	}
	r.log("Referenced row [%v] inserted into %s.", ref, grid)
	return nil
}

// function is available for mocking
var getInsertStatementForReferenceRow = func() string {
	return "INSERT INTO relationships (uuid, " +
		"revision, " +
		"created, " +
		"updated, " +
		"createdBy, " +
		"updatedBy, " +
		"enabled, " +
		"gridUuid, " +
		"text1, " +
		"text2, " +
		"text3, " +
		"text4, " +
		"text5) " +
		"VALUES ($1, " +
		"1, " +
		"NOW(), " +
		"NOW(), " +
		"$2, " +
		"$2, " +
		"true, " +
		"$3, " +
		"$4, " +
		"$5, " +
		"$6, " +
		"$7, " +
		"$8)"
}

func getInsertStatementParametersForRefereceRow(r apiRequestParameters, grid *model.Grid, ref gridReferencePost, rowUuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, utils.GetNewUUID())
	parameters = append(parameters, r.userUuid)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, ref.ColumnName)
	parameters = append(parameters, grid.Uuid)
	parameters = append(parameters, rowUuid)
	parameters = append(parameters, ref.ToGridUuid)
	parameters = append(parameters, ref.ToUuid)
	return parameters
}

func postGridSetOwnership(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row) error {
	if grid.Uuid == model.UuidGrids {
		r.trace("postGridSetOwnership(%s, %v)", grid, addedRows)
		for _, row := range addedRows {
			ref := gridReferencePost{
				ColumnName: "relationship3",
				FromUuid:   row.Uuid,
				ToGridUuid: model.UuidUsers,
				ToUuid:     r.userUuid,
			}
			err := postInsertReferenceRow(r, grid, addedRows, ref)
			if err != nil {
				_ = r.rollbackTransaction()
				return err
			}
		}
	}
	return nil
}

func getUuidFromRowsForTmpUuid(r apiRequestParameters, addedRows []*model.Row, tmpUuid string) string {
	for _, row := range addedRows {
		r.trace("getUuidFromRowsForTmpUuid() - row.TmpUuid=%v, row.Uuid=%v", row.TmpUuid, row.Uuid)
		if row.TmpUuid == tmpUuid {
			return row.Uuid
		}
	}
	return tmpUuid
}

func postDeleteReferenceRow(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost) error {
	query := getDeleteReferenceRowStatement()
	parms := getDeleteReferenceRowStatementParameters(r, grid, ref)
	r.trace("postDeleteReferenceRow(%s, %v, %v) query=%s ; params=%s", grid, addedRows, ref, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Delete referenced row error: %v.", err)
	}
	if grid.Uuid == model.UuidGrids {
		removeGridFromCache(ref.FromUuid)
	}
	r.log("Referenced row [%v] deleted.", ref)
	return nil
}

// function is available for mocking
var getDeleteReferenceRowStatement = func() string {
	return "UPDATE relationships SET enabled = false WHERE gridUuid = $1 AND text1 = $2 AND text2 = $3 AND text3 = $4 AND text4 = $5 AND text5 = $6"
}

func getDeleteReferenceRowStatementParameters(r apiRequestParameters, grid *model.Grid, ref gridReferencePost) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, ref.ColumnName)
	parameters = append(parameters, grid.Uuid)
	parameters = append(parameters, ref.FromUuid)
	parameters = append(parameters, ref.ToGridUuid)
	parameters = append(parameters, ref.ToUuid)
	return parameters
}
