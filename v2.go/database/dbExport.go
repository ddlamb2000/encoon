// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"
	"os"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"gopkg.in/yaml.v2"
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
	out, err := convertYaml(rowSet)
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
var convertYaml = func(rowSet []model.Row) ([]byte, error) {
	return yaml.Marshal(rowSet)
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
		"text1, " +
		"text2, " +
		"text3, " +
		"text4, " +
		"text5, " +
		"text6, " +
		"text7, " +
		"text8, " +
		"text9, " +
		"text10, " +
		"int1, " +
		"int2, " +
		"int3, " +
		"int4, " +
		"int5, " +
		"int6, " +
		"int7, " +
		"int8, " +
		"int9, " +
		"int10, " +
		"enabled, " +
		"revision "
}

// function is available for mocking
var getRowsQueryOutputForExportDb = func(row *model.Row) []any {
	return row.GetRowsQueryOutput()
}
