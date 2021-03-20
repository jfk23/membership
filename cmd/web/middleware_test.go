package main

import (
	"net/http"
	"testing"
)

func TestSulf(t *testing.T) {
	var mh myHandler
	r := NoSulf(&mh)

	switch r.(type) {
	case http.Handler :
		//do nothing
	default :
		t.Error("There is error with NoSulf() function.")
	}
}

func TestLoadSession(t *testing.T) {
	var mh myHandler
	r := LoadSession(&mh)

	switch r.(type) {
	case http.Handler :
		//do nothing
	default :
		t.Error("There is error with LoadSession() function.")
	}
}