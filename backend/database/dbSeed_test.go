// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/configuration"
)

func TestSeedDbIncorrectDb(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := SeedDb(context.Background(), "xxx", "/tmp/nothing.yml")
	if err == nil {
		t.Errorf("Can import database while it shouldn't: %v.", err)
	}
}

func TestSeedDbIncorrectFile(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	err := SeedDb(context.Background(), "test", "/tmp/nothing.yml")
	if err == nil {
		t.Errorf("Can import database while it shouldn't: %v.", err)
	}
}
