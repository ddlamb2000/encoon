// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"
	"strings"
	"testing"

	"d.lambert.fr/encoon/configuration"
	_ "github.com/lib/pq"
)

func TestGetRowsColumnDefinitions(t *testing.T) {
	got := getRowsColumnDefinitions()
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
	err = RecreateDb(context.Background(), db, "xxx")
	expect := "Only test database can be recreated."
	if err.Error() != expect {
		t.Errorf(`expect %v and found %v.`, expect, err)
	}
}

func TestRecreateDb3(t *testing.T) {
	getRowsColumnDefinitionsImpl := getRowsColumnDefinitions
	getRowsColumnDefinitions = func() string { return "x x x" } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := RecreateDb(context.Background(), db, dbName); err == nil {
		t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
	}
	getRowsColumnDefinitions = getRowsColumnDefinitionsImpl
}

func TestRecreateDb4(t *testing.T) {
	getMigrationInsertStatementImpl := getMigrationInsertStatement
	getMigrationInsertStatement = func() string { return "x x x" } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := RecreateDb(context.Background(), db, dbName); err == nil {
		t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
	}
	getMigrationInsertStatement = getMigrationInsertStatementImpl
}

func TestRecreateDb5(t *testing.T) {
	getMigrationStepsImpl := getMigrationSteps
	getMigrationSteps = func(string) map[int]string { return map[int]string{1: "x x x"} } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := RecreateDb(context.Background(), db, dbName); err == nil {
		t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
	}
	getMigrationSteps = getMigrationStepsImpl
}

func TestRecreateDb2(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := RecreateDb(context.Background(), db, dbName); err != nil {
		t.Errorf(`Can't recreate database %q: %v.`, dbName, err)
	}
}

func TestRecreateDb6(t *testing.T) {
	getDeletionStepsImpl := getDeletionSteps
	getDeletionSteps = func() map[int]string { return map[int]string{1: "x x x"} } // mock function
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := RecreateDb(context.Background(), db, dbName); err == nil {
		t.Errorf(`Can recreate database %q while it shouldn't be.`, dbName)
	}
	getDeletionSteps = getDeletionStepsImpl
}
