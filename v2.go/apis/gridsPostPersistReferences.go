// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type persistGridReferenceDataFunc func(context.Context, string, *sql.DB, string, string, *model.Grid, []*model.Row, gridReferencePost, string) error

func persistGridReferenceData(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, addedRows []*model.Row, refs []gridReferencePost, trace string, f persistGridReferenceDataFunc) error {
	for _, ref := range refs {
		err := f(ctx, dbName, db, userUuid, user, grid, addedRows, ref, trace)
		if err != nil {
			_ = RollbackTransaction(ctx, dbName, db, userUuid, user, trace)
			return err
		}
	}
	return nil
}

func postInsertReferenceRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost, trace string) error {
	utils.Trace(trace, "postInsertReferenceRow()")
	insertStatement := getInsertStatementForRefereceRow()
	utils.Trace(trace, "postInsertReferenceRow() - ref.FromUuid=%v, addedRows=%v", ref.FromUuid, addedRows)
	rowUuid := getUuidFromRowsForTmpUuid(addedRows, ref.FromUuid, trace)
	utils.Trace(trace, "postInsertReferenceRow() - rowUuid=%v", rowUuid)
	_, err := db.ExecContext(ctx, insertStatement, utils.GetNewUUID(), userUuid, model.UuidRelationships, ref.ColumnName, grid.Uuid, rowUuid, ref.ToGridUuid, ref.ToUuid)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Insert referenced row error on %q: %v.", dbName, user, insertStatement, err)
	}
	utils.Log("[%s] [%s] Referenced row [%v] inserted into %q.", dbName, user, ref, grid.GetUri())
	return err
}

func getUuidFromRowsForTmpUuid(addedRows []*model.Row, tmpUuid string, trace string) string {
	for _, row := range addedRows {
		utils.Trace(trace, "getUuidFromRowsForTmpUuid() - row.TmpUuid=%v, row.Uuid=%v", row.TmpUuid, row.Uuid)
		if row.TmpUuid == tmpUuid {
			return row.Uuid
		}
	}
	return tmpUuid
}

func getInsertStatementForRefereceRow() string {
	return "INSERT INTO rows (uuid, " +
		"version, " +
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

func postDeleteReferenceRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost, trace string) error {
	utils.Trace(trace, "postDeleteReferenceRow()")
	deleteStatement := getDeleteReferenceRowStatement()
	_, err := db.ExecContext(ctx, deleteStatement, model.UuidRelationships, ref.ColumnName, grid.Uuid, ref.FromUuid, ref.ToGridUuid, ref.ToUuid)
	if err != nil {
		return utils.LogAndReturnError("[%s] [%s] Delete referenced row error on %q: %v.", dbName, user, deleteStatement, err)
	}
	utils.Log("[%s] [%s] Referenced row [%v] delete into %q.", dbName, user, ref, grid.GetUri())
	return err
}

func getDeleteReferenceRowStatement() string {
	return "DELETE FROM rows WHERE gridUuid = $1 AND text1 = $2 AND text2 = $3 AND text3 = $4 AND text4 = $5 AND text5 = $6"
}
