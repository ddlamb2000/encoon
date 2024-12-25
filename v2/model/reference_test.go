// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package model

import (
	"testing"
)

func TestGetNewReference(t *testing.T) {
	reference := GetNewReference()
	if reference == nil {
		t.Errorf(`Isse when creating reference.`)
	}
}
