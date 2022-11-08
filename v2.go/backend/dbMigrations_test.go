// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/utils"
	_ "github.com/lib/pq"
)

func TestRecreateDb(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	if err := ConnectDbServers(utils.DatabaseConfigurations); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
	db := getDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := recreateDb(ctx, db, "xxx")
	expect := "[xxx] Only test database can be recreated."
	if err.Error() != expect {
		t.Errorf(`expect %v and found %v.`, expect, err)
	}
}
