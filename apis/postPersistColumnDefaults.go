// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"fmt"

	"d.lambert.fr/encoon/model"
)

func persistUpdateColumnDefaults(r apiRequest, grid *model.Grid, payload gridPost) error {
	gridUuids, _ := getGridsToUpdateWithColumnDefaults(r, grid, payload)
	r.log("persistUpdateColumnDefaults - gridUuids=%v", gridUuids)
	for _, gridUuid := range gridUuids {
		grid, err := getGridInstanceWithColumnsForGridsApi(r, gridUuid)
		if err != nil || grid == nil {
			return err
		}
		err = setGridsColumnDefaults(r, grid)
		if err != nil {
			return err
		}
	}
	return nil
}

func getGridsToUpdateWithColumnDefaults(r apiRequest, grid *model.Grid, payload gridPost) ([]string, error) {
	gridUuids := make([]string, 0)
	var mapGridUuids = make(map[string]bool)
	if grid.Uuid == model.UuidGrids {
		gridUuids = append(gridUuids, grid.Uuid)
	} else if grid.Uuid == model.UuidColumns {
		allRows := make([]*model.Row, 0)
		allRows = append(allRows, payload.RowsAdded...)
		allRows = append(allRows, payload.RowsEdited...)
		allRows = append(allRows, payload.RowsDeleted...)
		columnUuids := make([]string, 0)
		var mapColumnUuids = make(map[string]bool)
		for _, row := range allRows {
			if !mapColumnUuids[row.Uuid] {
				columnUuids = append(columnUuids, row.Uuid)
				mapColumnUuids[row.Uuid] = true
			}
		}
		for _, columnUuid := range columnUuids {
			r.trace("getGridsToUpdateWithColumnDefaults - columnUuid=%v", columnUuid)
			gridUuid, err := getGridUuidAttachedToColumn(r, columnUuid)
			if err != nil {
				return nil, err
			}
			if gridUuid != "" && !mapGridUuids[gridUuid] {
				gridUuids = append(gridUuids, gridUuid)
				mapGridUuids[gridUuid] = true
			}
		}
	}
	return gridUuids, nil
}

func setGridsColumnDefaults(r apiRequest, grid *model.Grid) error {
	var mapColumnIndexes = make(map[string]int64)
	var maxOrderNumber int64 = 0
	for _, column := range grid.Columns {
		prefix, index := column.GetColumnNamePrefixAndIndex()
		if prefix != "" && index > 0 && index > mapColumnIndexes[prefix] {
			mapColumnIndexes[prefix] = index
		}
		if column.OrderNumber != nil && *column.OrderNumber > maxOrderNumber {
			maxOrderNumber = *column.OrderNumber
		}
	}
	r.log("setGridsColumnDefaults(%v) - maxOrderNumber=%d, mapColumnIndexes = %v", grid, maxOrderNumber, mapColumnIndexes)
	for _, column := range grid.Columns {
		prefix, _ := column.GetColumnNamePrefixAndIndex()
		expectedPrefix := column.GetColumnNamePrefixFromType()
		r.trace("setGridsColumnDefaults - column %q ; prefix %q ; expectedPrefix %q", column.Label, prefix, expectedPrefix)
		if column.OrderNumber == nil || *column.OrderNumber == (int64)(0) {
			r.log("setGridsColumnDefaults - set column %q with order number %d", column.Label, maxOrderNumber+1)
			err := updateColumnOrderNumber(r, column.Uuid, maxOrderNumber+1)
			if err != nil {
				return err
			}
			maxOrderNumber += 1
		}
		if expectedPrefix != "" && (column.Name == "" || prefix != expectedPrefix) {
			newIndex := mapColumnIndexes[expectedPrefix] + 1
			newName := fmt.Sprintf("%s%d", expectedPrefix, newIndex)
			r.log("setGridsColumnDefaults - set column %q with name %s", column.Label, newName)
			err := updateColumnName(r, column.Uuid, newName)
			if err != nil {
				return err
			}
			mapColumnIndexes[expectedPrefix] = newIndex
		}
	}
	return nil
}

func updateColumnOrderNumber(r apiRequest, columnUuid string, orderNumber int64) error {
	query := getUpdateColumnOrderNumberQuery()
	r.trace("updateColumnOrderNumber(%s, %d) - query=%s", columnUuid, orderNumber, query)
	if err := r.execContext(query, model.UuidColumns, columnUuid, orderNumber, r.p.userUuid); err != nil {
		return r.logAndReturnError("Delete row error: %v.", err)
	}
	return nil
}

var getUpdateColumnOrderNumberQuery = func() string {
	return "UPDATE columns " +
		"SET int1 = $3, " +
		"updated = NOW(), " +
		"updatedBy = $4 " +
		"WHERE gridUuid = $1 " +
		"AND uuid = $2"
}

func updateColumnName(r apiRequest, columnUuid string, name string) error {
	query := getUpdateColumnNameQuery()
	r.trace("updateColumnName(%s, %s) - query=%s", columnUuid, name, query)
	if err := r.execContext(query, model.UuidColumns, columnUuid, name, r.p.userUuid); err != nil {
		return r.logAndReturnError("Delete row error: %v.", err)
	}
	return nil
}

var getUpdateColumnNameQuery = func() string {
	return "UPDATE columns " +
		"SET text2 = $3, " +
		"updated = NOW(), " +
		"updatedBy = $4 " +
		"WHERE gridUuid = $1 " +
		"AND uuid = $2"
}
