// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"strings"
	"testing"

	"d.lambert.fr/encoon/configuration"
	_ "github.com/lib/pq"
)

func TestConnectDbServers1(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	_, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func TestConnectDbServers2(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	_, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func TestConnectDbServers3(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	_, err := GetDbByName(dbName)
	if err != nil {
		t.Errorf(`Database %q not found.`, dbName)
	}
}

func TestConnectDbServers4(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "baddb"
	_, err := GetDbByName(dbName)
	got := err.Error()
	expect := "Unable to connect to database: dial tcp"
	if err == nil || !strings.Contains(got, expect) {
		t.Errorf("Got err %v instead of %v.", err, expect)
	}
}

func TestConnectDbServer(t *testing.T) {
	var conf configuration.DatabaseConfiguration
	conf.Host = "xxx"
	if _, err := connectDbServer(&conf, "root"); err == nil {
		t.Errorf(`Can connect to database?: %v.`, err)
	}
}

func TestGetDbByName(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	_, err := GetDbByName("")
	expect := "Missing database name parameter."
	if err == nil || err.Error() != expect {
		t.Errorf("Got err %v instead of %v.", err, expect)
	}
}
