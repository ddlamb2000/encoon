// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package core

import (
	"testing"
)

func TestGetUserByID1(t *testing.T) {
	LoadUsers()
	user, found := GetUserByID("c788a76d-4aa6-4073-8904-35a9b99a3289")
	if !found {
		t.Fatalf(`No user found %q`, user.ID.Uuid)
	}
	if user.ID.Uri != "root" {
		t.Fatalf(`Incorrect uri %q, expected %q`, user.ID.Uri, "root")
	}
}

func TestGetUserByID2(t *testing.T) {
	LoadUsers()
	user, found := GetUserByID("c788a76d-4aa6-4073-8904-35a9b99a3288")
	if found {
		t.Fatalf(`User found: %q`, user.ID.Uuid)
	}
}
