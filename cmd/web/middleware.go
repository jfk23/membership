package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func WriteToConsole(hd http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("page got hit!")
		hd.ServeHTTP(rw, r)
	})
}

func NoSulf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   appConfig.InProduction,
		Path:     "/",
	},)

	return csrfHandler
}

func LoadSession(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}
