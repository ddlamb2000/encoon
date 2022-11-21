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
		t.Errorf("Can't export daabase: %v.", err)
	}
}
