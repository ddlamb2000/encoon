// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

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

func TestExportDb4(t *testing.T) {
	getRowsQueryForExportDbImpl := getRowsQueryForExportDb
	getRowsQueryForExportDb = func() string { return "xxx" } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	getRowsQueryForExportDb = getRowsQueryForExportDbImpl
}

func TestExportDb5(t *testing.T) {
	getRowsQueryOutputForExportDbImpl := getRowsQueryOutputForExportDb
	getRowsQueryOutputForExportDb = func(*model.Row) []any { return nil } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	getRowsQueryOutputForExportDb = getRowsQueryOutputForExportDbImpl
}

func TestExportDb6(t *testing.T) {
	convertYamlImpl := convertYaml
	convertYaml = func(rowSet []model.Row) ([]byte, error) { return nil, errors.New("xxx") } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	convertYaml = convertYamlImpl
}

func TestExportDb7(t *testing.T) {
	exportToFileImpl := exportToFile
	exportToFile = func(f *os.File, out []byte) error { return errors.New("xxx") } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := ExportDb(context.Background(), "test", "/tmp/exportTestDb.yml")
	if err == nil {
		t.Errorf("Can export database while it shouldn't: %v.", err)
	}
	exportToFile = exportToFileImpl
}
