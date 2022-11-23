// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"

	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func postGridsRows(ct context.Context, dbName, userUuid, userName, gridUuid, uuid string, payload gridPost) (*model.Grid, []model.Row, int, bool, error) {
	r, cancel, err := createContextAndApiRequestParameters(ct, dbName, userUuid, userName)
	defer cancel()
	if err != nil {
		return nil, nil, 0, false, err
	}
	go func() {
		r.trace("postGridsRows()")
		if err := database.TestSleep(r.ctx, dbName, userName, r.db); err != nil {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Sleep interrupted: %v.", err)}
			return
		}
		grid, err := getGridForGridsApi(r, gridUuid)
		if err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := BeginTransaction(r); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsAdded, postInsertGridRow); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsEdited, postUpdateGridRow); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsDeleted, postDeleteGridRow); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesAdded, postInsertReferenceRow); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesRemoved, postDeleteReferenceRow); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := postGridSetOwnership(r, grid, payload.RowsAdded); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		if err := CommitTransaction(r); err != nil {
			r.ctxChan <- apiResponse{err: err}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, uuid, grid, true)
		if uuid != "" && rowSetCount == 0 {
			r.ctxChan <- apiResponse{grid: grid, err: r.logAndReturnError("Data not found.")}
			return
		}
		r.ctxChan <- apiResponse{grid: grid, rows: rowSet, rowCount: rowSetCount}
		r.trace("postGridsRows() - Done")
	}()
	select {
	case <-r.ctx.Done():
		r.trace("postGridsRows() - Cancelled")
		_ = RollbackTransaction(r)
		return nil, nil, 0, true, r.logAndReturnError("Post request has been cancelled: %v.", r.ctx.Err())
	case response := <-r.ctxChan:
		r.trace("postGridsRows() - OK")
		return response.grid, response.rows, response.rowCount, false, response.err
	}
}

func BeginTransaction(r apiRequestParameters) error {
	r.trace("beginTransaction()")
	_, err := r.db.ExecContext(r.ctx, "BEGIN")
	if err != nil {
		return r.logAndReturnError("Begin transaction error: %v.", err)
	}
	r.log("Begin transaction.")
	return err
}

func CommitTransaction(r apiRequestParameters) error {
	r.trace("commitTransaction()")
	_, err := r.db.ExecContext(r.ctx, "COMMIT")
	if err != nil {
		return r.logAndReturnError("Commit transaction error: %v.", err)
	}
	r.log("Commit transaction.")
	return err
}

func RollbackTransaction(r apiRequestParameters) error {
	r.trace("rollbackTransaction()")
	_, err := r.db.ExecContext(r.ctx, "ROLLBACK")
	if err != nil {
		return r.logAndReturnError("Rollback transaction error: %v.", err)
	}
	r.log("ROLLBACK transaction.")
	return err
}
