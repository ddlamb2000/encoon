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
	expect := "text01 text, text02 text, text03 text"
	if !strings.Contains(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
	expect2 := "int01 integer, int02 integer, int03 integer"
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
	expect := "[xxx] Only test database can be recreated."
	if err.Error() != expect {
		t.Errorf(`expect %v and found %v.`, expect, err)
	}
}
