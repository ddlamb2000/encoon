// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/configuration"
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
