// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/utils"
	_ "github.com/lib/pq"
)

func TestConnectDbServers(t *testing.T) {
	utils.LoadConfiguration("../", "configuration.yml")
	if err := ConnectDbServers(utils.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
	db := getDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := pingDb(ctx, db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
	if err := DisconnectDbServers(utils.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't disconnect to databases: %v.`, err)
	}
}

func TestConnectDbServer(t *testing.T) {
	var conf utils.Database
	conf.Host = "xxx"
	if err := connectDbServer(&conf); err == nil {
		t.Errorf(`Can connect to database?: %v.`, err)
	}
}
