package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/model"
)

var testApp config.AppConfig
var sessionManager *scs.SessionManager

func TestMain(m *testing.M) {

	gob.Register(model.Reservation{})

	testApp.InProduction = false

	testApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.Persist = true

	testApp.Session = sessionManager

	app = &testApp

	os.Exit(m.Run())
}

type MyWriter struct{}

func (w *MyWriter) Header() http.Header {
	var h http.Header
	return h
}

func (w *MyWriter) WriteHeader(code int) {
}

func (w *MyWriter) Write(b []byte) (int, error) {
	length := len(b)

	return length, nil
}