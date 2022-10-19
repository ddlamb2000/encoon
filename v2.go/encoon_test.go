// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package main

import "testing"

func TestSetAndStartServerHtml(t *testing.T) {
	srv := setAndStartServerHtml()
	if srv.Addr != portHtml {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}

func TestSetAndStartServerApi(t *testing.T) {
	srv := setAndStartServerApi()
	if srv.Addr != portApi {
		t.Fatalf(`Incorrect address %q`, srv.Addr)
	}
}

func TestInitWithLog(t *testing.T) {
	initWithLog()
}

func TestInitServers(t *testing.T) {
	srvHtml, srvApi := initServers()
	if srvHtml.Addr != portHtml {
		t.Fatalf(`Incorrect address %q`, srvHtml.Addr)
	}
	if srvApi.Addr != portApi {
		t.Fatalf(`Incorrect address %q`, srvApi.Addr)
	}
}
