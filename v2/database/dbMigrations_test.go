// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"strings"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	_ "github.com/lib/pq"
)

func TestGetRowsColumnDefinitions(t *testing.T) {
	got := getRowsColumnDefinitions(model.GetNewGrid(""))
	expect := "text1 text, text2 text, text3 text"
	if !strings.Contains(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
	expect2 := "int1 integer, int2 integer, int3 integer"
	if !strings.Contains(got, expect2) {
		t.Errorf(`Got %v instead of %v.`, got, expect2)
	}
}

func TestRecreateDb(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}

	t.Run("RecreateDb1", func(t *testing.T) {
		err = RecreateDb(context.Background(), db, "xxx")
		expect := "Only test database can be recreated."
		if err.Error() != expect {
			t.Errorf(`expect %v and found %v.`, expect, err)
		}
	})

	t.Run("RecreateDb2", func(t *testing.T) {
		getRowsColumnDefinitionsImpl := getRowsColumnDefinitions
		getRowsColumnDefinitions = func(*model.Grid) string { return "x x x" } // mock function
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
		}
		getRowsColumnDefinitions = getRowsColumnDefinitionsImpl
	})

	t.Run("RecreateDb3", func(t *testing.T) {
		getMigrationInsertStatementImpl := getMigrationInsertStatement
		getMigrationInsertStatement = func() string { return "x x x" } // mock function
		if err := migrateDbCommand(context.Background(), db, 0, 1, "SELECT pg_sleep(0.1)", dbName); err == nil {
			t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
		}
		getMigrationInsertStatement = getMigrationInsertStatementImpl
	})

	t.Run("RecreateDb4", func(t *testing.T) {
		getMigrationStepsImpl := getMigrationSteps
		getMigrationSteps = func(string) map[int]string { return map[int]string{1: "x x x"} } // mock function
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
		}
		getMigrationSteps = getMigrationStepsImpl
	})

	t.Run("RecreateDb5", func(t *testing.T) {
		if err := RecreateDb(context.Background(), db, dbName); err != nil {
			t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
		}
	})

	t.Run("RecreateDb6", func(t *testing.T) {
		getDeletionStepsImpl := getDeletionSteps
		getDeletionSteps = func() map[int]string { return map[int]string{1: "x x x"} } // mock function
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
		}
		getDeletionSteps = getDeletionStepsImpl
	})

	t.Run("MigrateDb", func(t *testing.T) {
		if err := migrateDb(context.Background(), db, dbName); err != nil {
			t.Errorf(`Error for migrating an up-to-date database: %v.`, err)
		}
	})

	t.Run("seedDb1", func(t *testing.T) {
		configuration.LoadConfiguration("../testData/invalidConfiguration1.yml")
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Database %q found while it shouldn't be.`, dbName)
		}
	})

	t.Run("seedDb2", func(t *testing.T) {
		configuration.LoadConfiguration("../testData/invalidConfiguration2.yml")
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Database %q found while it shouldn't be.`, dbName)
		}
	})

	t.Run("seedDbDefect1", func(t *testing.T) {
		GetGridRowsQueryForSeedDataImpl := GetGridRowsQueryForSeedData
		GetGridRowsQueryForSeedData = func(grid *model.Grid) string { return "xxx" }
		configuration.LoadConfiguration("../testData/validConfiguration1.yml")
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Database %q found while it shouldn't be.`, dbName)
		}
		GetGridRowsQueryForSeedData = GetGridRowsQueryForSeedDataImpl
	})

	t.Run("seedDbDefect2", func(t *testing.T) {
		GetGridInsertStatementForSeedRowDbImpl := GetGridInsertStatementForSeedRowDb
		GetGridInsertStatementForSeedRowDb = func(grid *model.Grid) string { return "xxx" }
		configuration.LoadConfiguration("../testData/validConfiguration1.yml")
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Database %q found while it shouldn't be.`, dbName)
		}
		GetGridInsertStatementForSeedRowDb = GetGridInsertStatementForSeedRowDbImpl
	})

	t.Run("seedDbDefect3", func(t *testing.T) {
		GetGridUpdateStatementForSeedRowDbImpl := GetGridUpdateStatementForSeedRowDb
		GetGridUpdateStatementForSeedRowDb = func(grid *model.Grid) string { return "xxx" }
		configuration.LoadConfiguration("../testData/validConfiguration1.yml")
		if err := RecreateDb(context.Background(), db, dbName); err == nil {
			t.Errorf(`Database %q found while it shouldn't be.`, dbName)
		}
		GetGridUpdateStatementForSeedRowDb = GetGridUpdateStatementForSeedRowDbImpl
	})

	t.Run("seedDb3", func(t *testing.T) {
		configuration.LoadConfiguration("../testData/validConfiguration1.yml")
		dbName := "test"
		db, err := GetDbByName(dbName)
		if err != nil {
			t.Errorf(`Database %q not found.`, dbName)
		}
		if err := RecreateDb(context.Background(), db, dbName); err != nil {
			t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
		}
		err = seedDb(context.Background(), db, dbName)
		if err != nil {
			t.Errorf(`Error when seeding database %q: %v`, dbName, err)
		}
	})
}
