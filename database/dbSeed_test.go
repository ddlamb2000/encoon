// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/configuration"
)

func TestSeedDb(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}

	t.Run("seedDb", func(t *testing.T) {
		if err := RecreateDb(context.Background(), db, dbName); err != nil {
			t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
		}
		err = seedDb(context.Background(), db, dbName)
		if err != nil {
			t.Errorf(`Error when seeding database %q: %v`, dbName, err)
		}
	})
}
