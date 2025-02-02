// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"context"

	"d.lambert.fr/encoon/database"
)

func PostGridsRows(ct context.Context, p ApiParameters, payload GridPost) GridResponse {
	r, cancel, err := createContextAndApiRequest(ct, p)
	defer cancel()
	t := r.startTiming()
	defer r.stopTiming("postGridsRows()", t)
	if err != nil {
		return GridResponse{Err: err}
	}
	go func() {
		r.trace("postGridsRows()")
		database.Sleep(r.ctx, p.DbName, p.UserName, r.db)
		grid, err := getGridForGridsApi(r, p.GridUuid)
		if err != nil {
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		} else if grid == nil {
			r.ctxChan <- GridResponse{Err: r.logAndReturnError("Data not found.")}
			return
		}
		canViewRows, canEditRows, canAddRows, canEditGrid := grid.GetViewEditAccessFlags(p.UserUuid)
		if !canViewRows {
			r.ctxChan <- GridResponse{Err: r.logAndReturnError("Access forbidden"), Forbidden: true}
			return
		}
		if err := r.beginTransaction(); err != nil {
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := postInsertTransaction(r); err != nil {
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsAdded, postInsertGridRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsEdited, postUpdateGridRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistGridRowData(r, grid, payload.RowsDeleted, postDeleteGridRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, defaultReferenceValues(r, payload), postInsertReferenceRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesAdded, postInsertReferenceRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistGridReferenceData(r, grid, payload.RowsAdded, payload.ReferenceValuesRemoved, postDeleteReferenceRow); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := persistUpdateColumnDefaults(r, grid, payload); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := postGridSetOwnership(r, grid, payload.RowsAdded); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		if err := r.commitTransaction(); err != nil {
			_ = r.rollbackTransaction()
			r.ctxChan <- GridResponse{Err: err, System: true}
			return
		}
		rowSet, rowSetCount, err := getRowSetForGridsApi(r, grid, p.Uuid, true, true)
		if p.Uuid != "" && rowSetCount == 0 {
			r.ctxChan <- GridResponse{Grid: grid, Err: r.logAndReturnError("Data not found.")}
			return
		}
		response := GridResponse{
			Err:                    err,
			Grid:                   grid,
			Rows:                   rowSet,
			CountRows:              rowSetCount,
			CanViewRows:            canViewRows,
			CanEditRows:            canEditRows,
			CanAddRows:             canAddRows,
			CanEditGrid:            canEditGrid,
			GridUuid:               p.GridUuid,
			ColumnUuid:             p.ColumnUuid,
			Uuid:                   p.Uuid,
			FilterColumnOwned:      p.FilterColumnOwned,
			FilterColumnName:       p.FilterColumnName,
			FilterColumnGridUuid:   p.FilterColumnGridUuid,
			FilterColumnValue:      p.FilterColumnValue,
			RowsAdded:              payload.RowsAdded,
			RowsEdited:             payload.RowsEdited,
			RowsDeleted:            payload.RowsDeleted,
			ReferenceValuesAdded:   payload.ReferenceValuesAdded,
			ReferenceValuesRemoved: payload.ReferenceValuesRemoved}
		if !grid.IsMetadata() {
			go postGridEmbeddings(ct, p, grid, response)
		}
		r.trace("postGridsRows() - Done")
		r.ctxChan <- response
	}()
	select {
	case <-r.ctx.Done():
		r.trace("postGridsRows() - Cancelled")
		_ = r.rollbackTransaction()
		return GridResponse{Err: r.logAndReturnError("Post request has been cancelled: %v.", r.ctx.Err()), TimeOut: true}
	case response := <-r.ctxChan:
		r.trace("postGridsRows() - OK")
		return response
	}
}
