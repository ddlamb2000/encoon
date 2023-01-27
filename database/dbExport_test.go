// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

func TestExportDb(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err != nil {
		t.Errorf("Can't export database: %v.", err)
	}
}

func TestExportDb2(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/xxx/yyy/zzz/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
}

func TestExportDb3(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "xxx", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
}

func TestExportDbDefect1(t *testing.T) {
	GetGridRowsQueryForExportDbImpl := GetGridRowsQueryForExportDb
	GetGridRowsQueryForExportDb = func(*model.Grid) string { return "xxx" }
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	GetGridRowsQueryForExportDb = GetGridRowsQueryForExportDbImpl
}

func TestExportDbDefect2(t *testing.T) {
	GetRowsQueryOutputImpl := GetRowsQueryOutput
	GetRowsQueryOutput = func(*model.Row) []any { return make([]any, 0) }
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	GetRowsQueryOutput = GetRowsQueryOutputImpl
}

func TestExportDbDefect3(t *testing.T) {
	convertJsonImpl := convertJson
	convertJson = func(rowSet []model.Row) ([]byte, error) { return nil, errors.New("xxx") } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	convertJson = convertJsonImpl
}

func TestExportDefect4(t *testing.T) {
	exportToFileImpl := exportToFile
	exportToFile = func(f *os.File, out []byte) error { return errors.New("xxx") } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	exportToFile = exportToFileImpl
}
