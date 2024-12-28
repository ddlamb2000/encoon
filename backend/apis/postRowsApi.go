// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"context"

	"d.lambert.fr/encoon/database"
)

func postGridsRows(ct context.Context, uri string, p HtmlParameters, payload gridPost) ApiResponse {
	r, cancel, err := createContextAndApiRequest(ct, p, uri)
	defer cancel()
	t := r.startTiming()
	defer r.stopTiming("postGridsRows()", t)
	if err != nil {
		return ApiResponse{Err: err}
	}
	go func() {
		r.trace("postGridsRows()")
		database.Sleep(r.ctx, p.DbName, p.userName, r.db)
		grid, err := getGridForGridsApi(r, p.GridUuid)
		if err != nil {
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		} else if grid == nil {
			r.ctxChan <- ApiResponse{Err: r.logAndReturnError("Data not found.")}
			return
		}
		canViewRows, canEditRows, canAddRows, canEditGrid := grid.GetViewEditAccessFlags(p.userUuid)
		if !canViewRows {
			r.ctxChan <- ApiResponse{Err: r.logAndReturnError("Access forbidden."), Forbidden: true}
			return
		}
		if err := r.beginTransaction(); err != nil {
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := postInsertTransaction(r); err != nil {
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsAdded, postInsertGridRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsEdited, postUpdateGridRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsDeleted, postDeleteGridRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, defaultReferenceValues(r, payload), postInsertReferenceRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesAdded, postInsertReferenceRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesRemoved, postDeleteReferenceRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := persistUpdateColumnDefaults(r, grid, payload); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := postGridSetOwnership(r, grid, payload.RowsAdded); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		if err := r.commitTransaction(); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- ApiResponse{Err: err, System: true}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, grid, p.Uuid, true, true)
		if p.Uuid != "" && rowSetCount == 0 {
			r.ctxChan <- ApiResponse{Grid: grid, Err: r.logAndReturnError("Data not found.")}
			return
		}
		r.ctxChan <- ApiResponse{
			Err:                    err,
			Grid:                   grid,
			Rows:                   rowSet,
			CountRows:              rowSetCount,
			CanViewRows:            canViewRows,
			CanEditRows:            canEditRows,
			CanAddRows:             canAddRows,
			CanEditGrid:            canEditGrid,
			RowsAdded:              payload.RowsAdded,
			RowsEdited:             payload.RowsEdited,
			RowsDeleted:            payload.RowsDeleted,
			ReferenceValuesAdded:   payload.ReferenceValuesAdded,
			ReferenceValuesRemoved: payload.ReferenceValuesRemoved}
		r.trace("postGridsRows() - Done")
	}()
	select {
	case <-r.ctx.Done():
		r.trace("postGridsRows() - Cancelled")
		_ = r.rollbackTransaction()
		return ApiResponse{Err: r.logAndReturnError("Post request has been cancelled: %v.", r.ctx.Err()), TimeOut: true}
	case response := <-r.ctxChan:
		r.trace("postGridsRows() - OK")
		return response
	}
}
