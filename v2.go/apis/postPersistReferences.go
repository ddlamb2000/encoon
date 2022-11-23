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
			_ = RollbackTransaction(r)
			return err
		}
	}
	return nil
}

func postInsertReferenceRow(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost) error {
	r.trace("postInsertReferenceRow()")
	insertStatement := getInsertStatementForRefereceRow()
	rowUuid := getUuidFromRowsForTmpUuid(r, addedRows, ref.FromUuid)
	_, err := r.db.ExecContext(r.ctx, insertStatement, utils.GetNewUUID(), r.userUuid, model.UuidRelationships, ref.ColumnName, grid.Uuid, rowUuid, ref.ToGridUuid, ref.ToUuid)
	if err != nil {
		return r.logAndReturnError("Insert referenced row error on %q: %v.", insertStatement, err)
	}
	r.log("Referenced row [%v] inserted into %q.", ref, grid.Uuid)
	return err
}

func postGridSetOwnership(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row) error {
	r.trace("postGridSetOwnership()")
	if grid.Uuid == model.UuidGrids {
		for _, row := range addedRows {
			ref := gridReferencePost{
				ColumnName: "relationship3",
				FromUuid:   row.Uuid,
				ToGridUuid: model.UuidUsers,
				ToUuid:     r.userUuid,
			}
			err := postInsertReferenceRow(r, grid, addedRows, ref)
			if err != nil {
				_ = RollbackTransaction(r)
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

func getInsertStatementForRefereceRow() string {
	return "INSERT INTO rows (uuid, " +
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

func postDeleteReferenceRow(r apiRequestParameters, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost) error {
	r.trace("postDeleteReferenceRow()")
	deleteStatement := getDeleteReferenceRowStatement()
	_, err := r.db.ExecContext(r.ctx, deleteStatement, model.UuidRelationships, ref.ColumnName, grid.Uuid, ref.FromUuid, ref.ToGridUuid, ref.ToUuid)
	if err != nil {
		return r.logAndReturnError("Delete referenced row error on %q: %v.", deleteStatement, err)
	}
	r.log("Referenced row [%v] delete into %q.", ref, grid.Uuid)
	return err
}

func getDeleteReferenceRowStatement() string {
	return "DELETE FROM rows WHERE gridUuid = $1 AND text1 = $2 AND text2 = $3 AND text3 = $4 AND text4 = $5 AND text5 = $6"
}
