// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package database

import (
	"testing"

	"d.lambert.fr/encoon/configuration"
	_ "github.com/lib/pq"
)

func TestIsDbAuthorized(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	ConnectDbServers(configuration.GetConfiguration().Databases)
	dbName := "test"
	uuid, firstName, lastName, _, err := IsDbAuthorized(dbName, "root", "dGVzdA==")
	if err != nil {
		t.Errorf("Can't authenticate: %v.", err)
	}
	if err == nil && (uuid == "" || firstName == "" || lastName == "") {
		t.Errorf("Can't retrieve identifiers: %v, %v, %v.", uuid, firstName, lastName)
	}
	uuid, firstName, lastName, _, err2 := IsDbAuthorized(dbName, "rot", "dGVzdA==")
	if err2 == nil {
		t.Errorf("Can authenticate with a wrong id: %v, %v, %v.", uuid, firstName, lastName)
	}
	uuid, firstName, lastName, _, err3 := IsDbAuthorized(dbName, "root", "========")
	if err3 == nil {
		t.Errorf("Can authenticate with a wrong password: %v, %v, %v.", uuid, firstName, lastName)
	}
	uuid, firstName, lastName, _, err4 := IsDbAuthorized("====", "root", "========")
	if err4 == nil {
		t.Errorf("Can authenticate on a dummy database: %v, %v, %v.", uuid, firstName, lastName)
	}
}
