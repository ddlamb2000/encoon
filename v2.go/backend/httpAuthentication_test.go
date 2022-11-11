// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
)

func TestGetNewToken(t *testing.T) {
	configuration.LoadConfiguration("../", "configuration.yml")
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().HttpServer.JwtExpiration) * time.Minute)
	token, err := getNewToken("test", "root", "0", "root", "root", expiration)
	if err != nil {
		t.Errorf("Token can't be created: %v.", err)
	}
	if len(token) < 20 {
		t.Errorf("Token doesn't seem to be a token: %v.", token)
	}
}
