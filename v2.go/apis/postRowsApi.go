// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"

	"d.lambert.fr/encoon/database"
)

func postGridsRows(ct context.Context, uri, dbName, userUuid, userName, gridUuid, uuid string, payload gridPost) apiResponse {
	r, cancel, err := createContextAndApiRequestParameters(ct, dbName, userUuid, userName, uri)
	defer cancel()
	if err != nil {
		return apiResponse{err: err}
	}
	go func() {
		r.trace("postGridsRows()")
		database.Sleep(r.ctx, dbName, userName, r.db)
		grid, err := getGridForGridsApi(r, gridUuid)
		if err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		} else if grid == nil {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Data not found.")}
			return
		}
		canViewRows, canEditRows, canEditOwnedRows, canAddRows := grid.GetViewEditAccessFlags(r.userUuid)
		if !canViewRows {
			r.ctxChan <- apiResponse{err: r.logAndReturnError("Access forbidden."), forbidden: true}
			return
		}
		if err := r.beginTransaction(); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := postInsertTransaction(r); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsAdded, postInsertGridRow); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsEdited, postUpdateGridRow); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsDeleted, postDeleteGridRow); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesAdded, postInsertReferenceRow); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesRemoved, postDeleteReferenceRow); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := postGridSetOwnership(r, grid, payload.RowsAdded); err != nil {
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		if err := r.commitTransaction(); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- apiResponse{err: err, system: true}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, grid, uuid, true, true)
		if uuid != "" && rowSetCount == 0 {
			r.ctxChan <- apiResponse{grid: grid, err: r.logAndReturnError("Data not found.")}
			return
		}
		r.ctxChan <- apiResponse{grid: grid, rows: rowSet, rowCount: rowSetCount, canViewRows: canViewRows, canEditRows: canEditRows, canEditOwnedRows: canEditOwnedRows, canAddRows: canAddRows}
		r.trace("postGridsRows() - Done")
	}()
	select {
	case <-r.ctx.Done():
		r.trace("postGridsRows() - Cancelled")
		_ = r.rollbackTransaction()
		return apiResponse{err: r.logAndReturnError("Post request has been cancelled: %v.", r.ctx.Err()), timeOut: true}
	case response := <-r.ctxChan:
		r.trace("postGridsRows() - OK")
		return response
	}
}
