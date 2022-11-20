// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"database/sql"

	"d.lambert.fr/encoon/model"
)

type persistGridReferenceDataFunc func(context.Context, string, *sql.DB, string, string, *model.Grid, *gridReferencePost, string) error

func persistGridReferenceData(ctx context.Context, dbName string, db *sql.DB, userUuid, user string, grid *model.Grid, rows []gridReferencePost, trace string, f persistGridReferenceDataFunc) error {
	for _, row := range rows {
		err := f(ctx, dbName, db, userUuid, user, grid, &row, trace)
		if err != nil {
			_ = RollbackTransaction(ctx, dbName, db, userUuid, user, trace)
			return err
		}
	}
	return nil
}
