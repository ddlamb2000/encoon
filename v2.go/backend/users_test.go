// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"testing"
)

func TestUserAsString(t *testing.T) {
	user := user{FirstName: "System", LastName: "Encoon"}
	out := user.String()
	expected := "System Encoon"
	if out != expected {
		t.Errorf("Incorrect user as string: found %q instead of %q.", out, expected)
	}
}

func TestGetUserByID1(t *testing.T) {
	loadUsers()
	user, found := GetUserByID("c788a76d-4aa6-4073-8904-35a9b99a3289")
	if !found {
		t.Errorf("No user found %q.", user.entity.Uuid)
	}
	if user.FirstName != "Root" {
		t.Errorf("Incorrect FirstName %q, expected %q.", user.FirstName, "root")
	}
}

func TestGetUserByID2(t *testing.T) {
	loadUsers()
	user, found := GetUserByID("c788a76d-4aa6-4073-8904-35a9b99a3288")
	if found {
		t.Errorf("User found: %q.", user.entity.Uuid)
	}
}
