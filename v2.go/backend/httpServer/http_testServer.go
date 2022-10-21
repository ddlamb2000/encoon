// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package httpServer

import "testing"

func TestSetAndStartServer(t *testing.T) {
	SetAndStartServer()
	if srv.Addr != httpPort {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}
