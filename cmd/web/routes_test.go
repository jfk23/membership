package main

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/jfk23/gobookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app *config.AppConfig

	r := Routes(app)

	switch r.(type) {
	case *chi.Mux:
		//cool do nothing
	default :
		t.Error("There is error with Route() function.")
	}

}