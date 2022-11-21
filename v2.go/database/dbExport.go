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
	configuration.Trace(dbName, "", "ExportDb()")

	rows, err := db.QueryContext(ct, getRowsQueryForExportDb())
	if err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error when querying rows: %v.", err)
	}
	defer rows.Close()
	var rowSet = make([]model.Row, 0)
	for rows.Next() {
		var row = new(model.Row)
		if err := rows.Scan(getRowsQueryOutputForExportDb(row)...); err != nil {
			return configuration.LogAndReturnError(dbName, "", "Error when exporting rows: %v.", err)
		}
		rowSet = append(rowSet, *row)
	}
	configuration.Trace(dbName, "", "ExportDb() - end of fetching rows.")
	out, err := yaml.Marshal(rowSet)
	if err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error when marshalling rows: %v.", err)
	}
	_, err = f.Write(out)
	if err != nil {
		return configuration.LogAndReturnError(dbName, "", "Error when writing file: %v.", err)
	}
	configuration.Trace(dbName, "", "ExportDb() - done.")
	return nil
}

func getRowsQueryForExportDb() string {
	selectStr := getRowsQueryColumnsForExportDb()
	fromStr := " FROM rows "
	orderByStr := " ORDER BY griduuid, uuid "
	return selectStr + fromStr + orderByStr
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
		"version "
}

func getRowsQueryOutputForExportDb(row *model.Row) []any {
	output := make([]any, 0)
	output = append(output, &row.Uuid)
	output = append(output, &row.GridUuid)
	output = append(output, &row.Created)
	output = append(output, &row.CreatedBy)
	output = append(output, &row.Updated)
	output = append(output, &row.UpdatedBy)
	output = append(output, &row.Text1)
	output = append(output, &row.Text2)
	output = append(output, &row.Text3)
	output = append(output, &row.Text4)
	output = append(output, &row.Text5)
	output = append(output, &row.Text6)
	output = append(output, &row.Text7)
	output = append(output, &row.Text8)
	output = append(output, &row.Text9)
	output = append(output, &row.Text10)
	output = append(output, &row.Int1)
	output = append(output, &row.Int2)
	output = append(output, &row.Int3)
	output = append(output, &row.Int4)
	output = append(output, &row.Int5)
	output = append(output, &row.Int6)
	output = append(output, &row.Int7)
	output = append(output, &row.Int8)
	output = append(output, &row.Int9)
	output = append(output, &row.Int10)
	output = append(output, &row.Enabled)
	output = append(output, &row.Version)
	return output
}
