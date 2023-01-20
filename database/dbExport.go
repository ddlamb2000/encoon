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
	flags := os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(exportFileName, flags, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	db, err := GetDbByName(dbName)
	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ct, getRowsQueryForExportDb())
	if err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error when querying rows: %v.", err)
	}
	defer rows.Close()
	rowSet := make([]model.Row, 0)
	for rows.Next() {
		row := model.GetNewRow()
		if err := rows.Scan(getRowsQueryOutputForExportDb(row)...); err != nil {
			return configuration.LogAndReturnError(dbName, "", "Error when exporting rows: %v.", err)
		}
		rowSet = append(rowSet, *row)
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
	return json.Marshal(rowSet)
}

// function is available for mocking
var exportToFile = func(f *os.File, out []byte) error {
	_, err := f.Write(out)
	return err
}

// function is available for mocking
var getRowsQueryForExportDb = func() string {
	return getRowsQueryColumnsForExportDb() + "FROM rows ORDER BY created"
}

func getRowsQueryColumnsForExportDb() string {
	return "SELECT uuid, " +
		"gridUuid, " +
		"created, " +
		"createdBy, " +
		"updated, " +
		"updatedBy, " +
		getRowsColumnDefinitions() +
		"enabled, " +
		"revision "
}

// function is available for mocking
var getRowsQueryOutputForExportDb = func(row *model.Row) []any {
	return row.GetRowsQueryOutput()
}
