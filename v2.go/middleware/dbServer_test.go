// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/utils"
)

func TestPing(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	ConnectDbServers(utils.DatabaseConfigurations)
	dbName := "test"
	db := getDbByName(dbName)
	if db == nil {
		t.Fatalf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := pingDb(ctx, db)
	if !err {
		t.Fatalf(`Database %q doesn't respond to ping.`, dbName)
	}
}
