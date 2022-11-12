// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/configuration"
	_ "github.com/lib/pq"
)

func TestConnectDbServers1(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	if err := ConnectDbServers(configuration.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
	db := GetDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func TestConnectDbServers2(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	if err := ConnectDbServers(configuration.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
	db := GetDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := PingDb(ctx, db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
}

func TestConnectDbServers3(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	if err := ConnectDbServers(configuration.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't connect to databases: %v.`, err)
	}
	dbName := "test"
	db := GetDbByName(dbName)
	if db == nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	if err := PingDb(ctx, db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
	if err := DisconnectDbServers(configuration.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't disconnect to databases: %v.`, err)
	}
	if err := ConnectDbServers(configuration.GetConfiguration().Databases); err != nil {
		t.Errorf(`Can't re-connect to databases: %v.`, err)
	}
}

func TestConnectDbServer(t *testing.T) {
	var conf configuration.DatabaseConfiguration
	conf.Host = "xxx"
	if err := connectDbServer(&conf); err == nil {
		t.Errorf(`Can connect to database?: %v.`, err)
	}
}
