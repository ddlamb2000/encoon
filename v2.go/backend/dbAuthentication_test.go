// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"testing"

	"d.lambert.fr/encoon/utils"
	_ "github.com/lib/pq"
)

func TestIsDbAuthorized(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	ConnectDbServers(utils.DatabaseConfigurations)
	dbName := "test"
	uuid, firstName, lastName, _, err := isDbAuthorized(dbName, "root", "dGVzdA==")
	if err != nil {
		t.Errorf("Can't authenticate: %v.", err)
	}
	if err == nil && (uuid == "" || firstName == "" || lastName == "") {
		t.Errorf("Can't retrieve identifiers: %v, %v, %v.", uuid, firstName, lastName)
	}
	uuid, firstName, lastName, _, err2 := isDbAuthorized(dbName, "rot", "dGVzdA==")
	if err2 == nil {
		t.Errorf("Can authenticate with a wrong id: %v, %v, %v.", uuid, firstName, lastName)
	}
	uuid, firstName, lastName, _, err3 := isDbAuthorized(dbName, "root", "========")
	if err3 == nil {
		t.Errorf("Can authenticate with a wrong password: %v, %v, %v.", uuid, firstName, lastName)
	}
	uuid, firstName, lastName, _, err4 := isDbAuthorized("====", "root", "========")
	if err4 == nil {
		t.Errorf("Can authenticate on a dummy database: %v, %v, %v.", uuid, firstName, lastName)
	}
}
