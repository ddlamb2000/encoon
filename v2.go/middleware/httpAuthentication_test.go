// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import (
	"testing"

	"d.lambert.fr/encoon/utils"
)

func TestGetNewToken(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	token, err := getNewToken("test", "root", "0", "root", "root")
	if err != nil {
		t.Fatalf("Token can't be created: %v.", err)
	}
	if len(token) < 20 {
		t.Fatalf("Token doesn't seem to be a token: %v.", token)
	}
}
