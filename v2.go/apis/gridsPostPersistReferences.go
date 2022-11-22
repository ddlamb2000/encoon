// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
)

type persistGridReferenceDataFunc func(context.Context, string, *sql.DB, string, string, *model.Grid, []*model.Row, gridReferencePost) error

func persistGridReferenceData(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, addedRows []*model.Row, refs []gridReferencePost, f persistGridReferenceDataFunc) error {
	for _, ref := range refs {
		err := f(ctx, dbName, db, userUuid, user, grid, addedRows, ref)
		if err != nil {
			_ = RollbackTransaction(ctx, dbName, db, userUuid, user)
			return err
		}
	}
	return nil
}

func postInsertReferenceRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost) error {
	configuration.Trace(dbName, user, "postInsertReferenceRow()")
	insertStatement := getInsertStatementForRefereceRow()
	configuration.Trace(dbName, user, "postInsertReferenceRow() - ref.FromUuid=%v, addedRows=%v", ref.FromUuid, addedRows)
	rowUuid := getUuidFromRowsForTmpUuid(addedRows, ref.FromUuid)
	configuration.Trace(dbName, user, "postInsertReferenceRow() - rowUuid=%v", rowUuid)
	_, err := db.ExecContext(ctx, insertStatement, utils.GetNewUUID(), userUuid, model.UuidRelationships, ref.ColumnName, grid.Uuid, rowUuid, ref.ToGridUuid, ref.ToUuid)
	if err != nil {
		return configuration.LogAndReturnError(dbName, user, "Insert referenced row error on %q: %v.", insertStatement, err)
	}
	configuration.Log(dbName, user, "Referenced row [%v] inserted into %q.", ref, grid.Uuid)
	return err
}

func getUuidFromRowsForTmpUuid(addedRows []*model.Row, tmpUuid string) string {
	for _, row := range addedRows {
		configuration.Trace("getUuidFromRowsForTmpUuid() - row.TmpUuid=%v, row.Uuid=%v", row.TmpUuid, row.Uuid)
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

func postDeleteReferenceRow(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, addedRows []*model.Row, ref gridReferencePost) error {
	configuration.Trace(dbName, user, "postDeleteReferenceRow()")
	deleteStatement := getDeleteReferenceRowStatement()
	_, err := db.ExecContext(ctx, deleteStatement, model.UuidRelationships, ref.ColumnName, grid.Uuid, ref.FromUuid, ref.ToGridUuid, ref.ToUuid)
	if err != nil {
		return configuration.LogAndReturnError(dbName, user, "Delete referenced row error on %q: %v.", deleteStatement, err)
	}
	configuration.Log(dbName, user, "Referenced row [%v] delete into %q.", ref, grid.Uuid)
	return err
}

func getDeleteReferenceRowStatement() string {
	return "DELETE FROM rows WHERE gridUuid = $1 AND text1 = $2 AND text2 = $3 AND text3 = $4 AND text4 = $5 AND text5 = $6"
}
