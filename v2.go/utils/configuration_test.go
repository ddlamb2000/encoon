// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package utils

import (
	"reflect"
	"testing"
)

func TestLoadMainConfiguration(t *testing.T) {
	path := "../"
	fileName := "configuration.yml"
	if err := loadMainConfiguration(path, fileName); err != nil {
		t.Errorf("Can't load configuration %q from path %q: %v.", fileName, path, err)
	}
}

func TestLoadMainConfiguration2(t *testing.T) {
	path := "../"
	fileName := "xxxx.yml"
	if err := loadMainConfiguration(path, fileName); err == nil {
		t.Errorf("Can load configuration %q from path %q.", fileName, path)
	}
}

func TestLoadMainConfiguration3(t *testing.T) {
	path := "../utils/"
	fileName := "configuration.go"
	if err := loadMainConfiguration(path, fileName); err == nil {
		t.Errorf("Can load configuration %q from path %q.", fileName, path)
	}
}

func TestGetRootAndPassword(t *testing.T) {
	root, password := GetRootAndPassword("xxx")
	if root != "" || password != "" {
		t.Errorf("Root or password isn't correct for database %q: %q and %q.", "xxx", root, password)
	}

}

func TestLoadConfiguration(t *testing.T) {
	dir := "../"
	fileName := "configuration.yml"
	if err := LoadConfiguration(dir, fileName); err != nil {
		t.Errorf("Can't load configurations from directory %q: %v.", dir, err)
	}

	dbName := "test"
	if !IsDatabaseEnabled(dbName) {
		t.Errorf("Database %q isn't enabled.", dbName)
	}

	secret := GetJWTSecret(dbName)
	expect := []byte{
		116, 101, 115, 116, 36, 50, 97, 36, 48, 56, 36,
		100, 99, 110, 50, 50, 118, 82, 70, 73, 90, 109,
		121, 119, 118, 100, 89, 66, 70, 118, 53, 121, 79,
		82, 99, 79, 105, 71, 85, 46, 90, 116, 113, 66,
		57, 83, 49, 100, 84, 99, 120, 115, 112, 86, 108,
		122, 97, 101, 108, 109, 90, 85, 80, 97,
	}
	if !reflect.DeepEqual(secret, expect) {
		t.Errorf("Jwt secret is wrong: %v instead of %v.", secret, expect)
	}

	root, password := GetRootAndPassword(dbName)
	expectPassword := "$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"
	if root != "root" || password != expectPassword {
		t.Errorf("Root or password isn't correct for database %q: %q and %q.", dbName, root, password)
	}
}

func TestLoadConfiguration2(t *testing.T) {
	dir := "../utils/"
	fileName := "configuration.yml"
	if err := LoadConfiguration(dir, fileName); err == nil {
		t.Errorf("Can't load configurations from directory %q: %v.", dir, err)
	}
}

func TestLoadConfiguration3(t *testing.T) {
	secret := GetJWTSecret("xxx")
	if secret != nil {
		t.Errorf("Invalid Jwt secret: %q.", secret)
	}
}
