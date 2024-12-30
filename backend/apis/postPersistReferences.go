// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type persistGridReferenceDataFunc func(ApiRequest, *model.Grid, []*model.Row, GridReferencePost) error

func persistGridReferenceData(r ApiRequest, grid *model.Grid, addedRows []*model.Row, refs []GridReferencePost, f persistGridReferenceDataFunc) error {
	for _, ref := range refs {
		err := f(r, grid, addedRows, ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func postInsertReferenceRow(r ApiRequest, grid *model.Grid, addedRows []*model.Row, ref GridReferencePost) error {
	query := getInsertStatementForReferenceRow()
	rowUuid := getUuidFromRowsForTmpUuid(r, addedRows, ref.FromUuid)
	parms := getInsertStatementParametersForReferenceRow(r, grid, ref, rowUuid)
	r.trace("postInsertReferenceRow(%s, %v, %v) - query=%s ; parms=%s", grid, addedRows, ref, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Insert referenced row error: %v.", err)
	}
	if err := removeAssociatedGridFromCache(r, grid, rowUuid); err != nil {
		return r.logAndReturnError("Error when getting data for cache deletion: %v.", err)
	}
	if ref.ToGridUuid == model.UuidGrids {
		removeGridFromCache(ref.ToUuid)
	}
	if grid.Uuid == model.UuidGrids && ref.ToGridUuid == model.UuidColumns {
		if err := removeAssociatedGridNotOwnedColumnFromCache(r, grid, ref.ToUuid); err != nil {
			return r.logAndReturnError("Error when getting data for cache deletion: %v.", err)
		}
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

func getInsertStatementParametersForReferenceRow(r ApiRequest, grid *model.Grid, ref GridReferencePost, rowUuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, utils.GetNewUUID())
	parameters = append(parameters, r.p.UserUuid)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, ref.ColumnName)
	if ref.Owned {
		parameters = append(parameters, grid.Uuid)
		parameters = append(parameters, rowUuid)
		parameters = append(parameters, ref.ToGridUuid)
		parameters = append(parameters, ref.ToUuid)
	} else {
		parameters = append(parameters, ref.ToGridUuid)
		parameters = append(parameters, ref.ToUuid)
		parameters = append(parameters, grid.Uuid)
		parameters = append(parameters, rowUuid)
	}
	return parameters
}

func postGridSetOwnership(r ApiRequest, grid *model.Grid, addedRows []*model.Row) error {
	if grid.Uuid == model.UuidGrids {
		r.trace("postGridSetOwnership(%s, %v)", grid, addedRows)
		for _, row := range addedRows {
			ref := GridReferencePost{
				Owned:      true,
				ColumnName: "relationship3",
				FromUuid:   row.Uuid,
				ToGridUuid: model.UuidUsers,
				ToUuid:     r.p.UserUuid,
			}
			err := postInsertReferenceRow(r, grid, addedRows, ref)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func postDeleteReferenceRow(r ApiRequest, grid *model.Grid, addedRows []*model.Row, ref GridReferencePost) error {
	query := getDeleteReferenceRowStatement()
	rowUuid := getUuidFromRowsForTmpUuid(r, addedRows, ref.FromUuid)
	parms := getDeleteReferenceRowStatementParameters(r, grid, ref, rowUuid)
	r.trace("postDeleteReferenceRow(%s, %v, %v) query=%s ; params=%s", grid, addedRows, ref, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Delete referenced row error: %v.", err)
	}
	if err := removeAssociatedGridFromCache(r, grid, ref.FromUuid); err != nil {
		return r.logAndReturnError("Error when getting data for cache deletion: %v.", err)
	}
	if ref.ToGridUuid == model.UuidGrids {
		removeGridFromCache(ref.ToUuid)
	}
	r.log("Referenced row [%v] deleted.", ref)
	return nil
}

// function is available for mocking
var getDeleteReferenceRowStatement = func() string {
	return "UPDATE relationships " +
		"SET enabled = false, " +
		"updated = NOW(), " +
		"updatedBy = $1 " +
		"WHERE gridUuid = $2 " +
		"AND text1 = $3 " +
		"AND text2 = $4 " +
		"AND text3 = $5 " +
		"AND text4 = $6 " +
		"AND text5 = $7"
}

func getDeleteReferenceRowStatementParameters(r ApiRequest, grid *model.Grid, ref GridReferencePost, rowUuid string) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, r.p.UserUuid)
	parameters = append(parameters, model.UuidRelationships)
	parameters = append(parameters, ref.ColumnName)
	if ref.Owned {
		parameters = append(parameters, grid.Uuid)
		parameters = append(parameters, rowUuid)
		parameters = append(parameters, ref.ToGridUuid)
		parameters = append(parameters, ref.ToUuid)
	} else {
		parameters = append(parameters, ref.ToGridUuid)
		parameters = append(parameters, ref.ToUuid)
		parameters = append(parameters, grid.Uuid)
		parameters = append(parameters, rowUuid)
	}
	return parameters
}

func getUuidFromRowsForTmpUuid(r ApiRequest, addedRows []*model.Row, tmpUuid string) string {
	for _, row := range addedRows {
		r.trace("getUuidFromRowsForTmpUuid() - row.TmpUuid=%v, row.Uuid=%v", row.TmpUuid, row.Uuid)
		if row.TmpUuid == tmpUuid {
			return row.Uuid
		}
	}
	return tmpUuid
}

func defaultReferenceValues(r ApiRequest, payload GridPost) []GridReferencePost {
	if r.p.filterColumnName == "" || r.p.filterColumnGridUuid == "" || r.p.filterColumnValue == "" {
		return nil
	}
	defaults := make([]GridReferencePost, 0)
	for _, rowAdded := range payload.RowsAdded {
		var foundReference = false
		for _, refAdded := range payload.ReferenceValuesAdded {
			if refAdded.FromUuid == rowAdded.TmpUuid &&
				refAdded.ColumnName == r.p.filterColumnName &&
				refAdded.ToGridUuid == r.p.filterColumnGridUuid &&
				refAdded.ToUuid == r.p.filterColumnValue &&
				refAdded.Owned == r.p.filterColumnOwned {
				foundReference = true
			}
		}
		if !foundReference {
			referencePost := GridReferencePost{
				FromUuid:   rowAdded.TmpUuid,
				ColumnName: r.p.filterColumnName,
				ToGridUuid: r.p.filterColumnGridUuid,
				ToUuid:     r.p.filterColumnValue,
				Owned:      r.p.filterColumnOwned,
			}
			r.trace("defaultReferenceValues() - referencePost=%v", referencePost)
			defaults = append(defaults, referencePost)
		}
	}
	return defaults
}
