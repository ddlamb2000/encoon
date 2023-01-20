// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"encoding/json"
	"os"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

func ExportDb(ct context.Context, dbName, exportFileName string) error {
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	f, err := os.OpenFile(exportFileName, flags, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	db, err := GetDbByName(dbName)
	if err != nil {
		return err
	}

	rowSet := make([]model.Row, 0)
	gridsToExport := []string{model.UuidGrids, model.UuidColumns, model.UuidRelationships, ""}
	for _, gridUuid := range gridsToExport {
		grid := model.GetNewGrid(gridUuid)
		tableName := grid.GetTableName()
		configuration.Log(dbName, "", "Export from %s.", tableName)
		query := grid.GetRowsQueryForExportDb()
		configuration.Trace(dbName, "", "ExportDb() - query=%s", query)
		rows, err := db.QueryContext(ct, query)
		if err != nil {
			return configuration.LogAndReturnError(dbName, "", "Error when querying rows: %v.", err)
		}
		defer rows.Close()
		for rows.Next() {
			row := model.GetNewRow()
			row.GridUuid = gridUuid
			if err := rows.Scan(row.GetRowsQueryOutput()...); err != nil {
				return configuration.LogAndReturnError(dbName, "", "Error when exporting rows: %v.", err)
			}
			rowSet = append(rowSet, *row)
		}
	}

	configuration.Trace(dbName, "", "ExportDb() - end of fetching rows.")
	out, err := convertJson(rowSet)
	if err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error when marshalling rows: %v.", err)
	}
	if err = exportToFile(f, out); err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error when writing file: %v.", err)
	}
	configuration.Log(dbName, "", "Export of database (%d rows) into %s completed.", len(rowSet), exportFileName)
	return nil
}

// function is available for mocking
var convertJson = func(rowSet []model.Row) ([]byte, error) {
	return json.MarshalIndent(rowSet, "", " ")
}

// function is available for mocking
var exportToFile = func(f *os.File, out []byte) error {
	_, err := f.Write(out)
	return err
}
