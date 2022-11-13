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
	dbName := "test"
	_, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func TestConnectDbServers2(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := PingDb(context.Background(), db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
}

func TestConnectDbServers3(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	dbName := "test"
	db, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
	if err := PingDb(context.Background(), db); err != nil {
		t.Errorf(`Database %q doesn't respond to ping: %v.`, dbName, err)
	}
}

func TestConnectDbServer(t *testing.T) {
	var conf configuration.DatabaseConfiguration
	conf.Host = "xxx"
	if _, err := connectDbServer(&conf); err == nil {
		t.Errorf(`Can connect to database?: %v.`, err)
	}
}
