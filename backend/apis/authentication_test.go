// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"testing"
	"time"

	"d.lambert.fr/encoon/configuration"
)

func TestGetNewToken(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	expiration := time.Now().Add(time.Duration(configuration.GetConfiguration().JwtExpiration) * time.Minute)
	token, err := getNewToken("test", "root", "0", "root", "root", expiration)
	if err != nil {
		t.Errorf("Token can't be created: %v.", err)
	}
	if len(token) < 20 {
		t.Errorf("Token doesn't seem to be a token: %v.", token)
	}
}
