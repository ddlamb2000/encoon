// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import "testing"

func TestSetAndStartServer(t *testing.T) {
	srv := setAndStartServer()
	if srv.Addr != httpPort {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}

func TestInitWithLog(t *testing.T) {
	initWithLog()
}

func TestInitServers(t *testing.T) {
	srv := initServers()
	if srv.Addr != httpPort {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}
