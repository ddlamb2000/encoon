// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package database

import (
	"context"
	"testing"

	"d.lambert.fr/encoon/configuration"
	_ "github.com/lib/pq"
)

func TestIsDbAuthorized1(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	_, _, _, _, err := IsDbAuthorized(context.Background(), dbName, "root", "dGVzdA==")
	if err != nil {
		t.Errorf("Can't authenticate: %v.", err)
	}
}

func TestIsDbAuthorized2(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	uuid, firstName, lastName, _, err2 := IsDbAuthorized(context.Background(), dbName, "rot", "dGVzdA==")
	if err2 == nil {
		t.Errorf("Can authenticate with a wrong id: %v, %v, %v.", uuid, firstName, lastName)
	}
}

func TestIsDbAuthorized3(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	uuid, firstName, lastName, _, err3 := IsDbAuthorized(context.Background(), dbName, "root", "========")
	if err3 == nil {
		t.Errorf("Can authenticate with a wrong password: %v, %v, %v.", uuid, firstName, lastName)
	}
}

func TestIsDbAuthorized4(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "===="
	uuid, firstName, lastName, _, err4 := IsDbAuthorized(context.Background(), dbName, "root", "========")
	if err4 == nil {
		t.Errorf("Can authenticate on a dummy database: %v, %v, %v.", uuid, firstName, lastName)
	}
}

func TestIsDbAuthorizedTimeOut(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	ForceTestSleepTimeAndTimeOutThreshold("test", 500, 200)
	defer ForceTestSleepTimeAndTimeOutThreshold("test", 0, 200)
	dbName := "test"
	_, _, _, _, err := IsDbAuthorized(context.Background(), dbName, "root", "dGVzdA==")
	if err == nil {
		t.Errorf("Can authenticate: %v!", err)
	}
	expect := "Authentication request has been cancelled: context deadline exceeded."
	if err == nil || err.Error() != expect {
		t.Errorf("Got err %v instead of %v.", err, expect)
	}
}
