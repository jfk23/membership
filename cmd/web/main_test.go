package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	_, e := run()
	if e !=nil {
		t.Error("There is error with main run() function")
	}
}