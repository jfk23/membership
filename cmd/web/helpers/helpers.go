package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/jfk23/gobookings/internal/config"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a 
}

func ClientError(r http.ResponseWriter, status int) {
	app.InfoLog.Println("We got an error from client with status code of: ", status)
	http.Error(r, http.StatusText(status), status)
}

func SeverError (r http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println("We got an error from server: ", trace)
	
	http.Error(r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticate(r *http.Request) bool {
	exist := app.Session.Exists(r.Context(), "user_id")
	return exist
} 