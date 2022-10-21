// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package middleware

import "testing"

func TestSetAndStartHttpServer(t *testing.T) {
	SetAndStartHttpServer()
	if srv.Addr != httpPort {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}
