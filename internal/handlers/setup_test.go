package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jfk23/gobookings/cmd/web/helpers"
	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/model"
	"github.com/jfk23/gobookings/internal/render"
	"github.com/justinas/nosurf"
)

var appConfig config.AppConfig
var sessionManager *scs.SessionManager

var funcmap = template.FuncMap{
	"humanTime" : render.HumanTime,
	"formatTime" : render.FormatTime,
	"iterate" : render.Iterate,
}
var pathToTemplate = "./../../templates"

func TestMain(m *testing.M) {

	gob.Register(model.Reservation{})
	gob.Register(model.User{})
	gob.Register(model.Room{})
	gob.Register(model.Restriction{})
	gob.Register(map[string]int{})

	appConfig.InProduction =false

	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = appConfig.InProduction
	sessionManager.Cookie.Persist = true

	appConfig.Session = sessionManager

	mailChen := make(chan model.MailData)
	appConfig.MailChan = mailChen

	listenForMail()
	

	ts, err := CreateTestTemplate()
	if err != nil {
		log.Fatal(err)
	}
	appConfig.CachedTemplate = ts
	appConfig.UseCache = true

	repo := CreateNewTestRepo(&appConfig)
	
	SetHandler(repo)

	render.NewRenderer(&appConfig)
	helpers.NewHelpers(&appConfig)

	os.Exit(m.Run())

}

func listenForMail() {
	go func() {
		for {
			_ = <- appConfig.MailChan
		}

	}()
}

func getRoutes() http.Handler {


	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	//mux.Use(NoSulf)
	mux.Use(LoadSession)


	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-room", Repo.Generals)
	mux.Get("/majors-room", Repo.Majors)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	
	mux.Get("/contact", Repo.Contact)
	mux.Get("/search", Repo.Search)
	mux.Post("/search", Repo.PostSearch)
	mux.Post("/search-json", Repo.PostSearchJson)

	mux.Get("/user/login", Repo.ShowLogin)
	mux.Post("/user/login", Repo.PostShowLogin)
	mux.Get("/user/logout", Repo.ShowLogout)

	mux.Get("/admin/dashboard", Repo.AdminDashBoard)
		
	mux.Get("/admin/reservations-new", Repo.AdminNewReservations)
	mux.Get("/admin/reservations-all", Repo.AdminAllReservations)
	mux.Get("/admin/reservations-calendar", Repo.AdminReservationsCalendar)
	mux.Post("/admin/reservations-calendar", Repo.AdminPostReservationsCalendar)
	mux.Get("/admin/reservation/{src}/{id}/show", Repo.AdminShowReservation)
	mux.Get("/admin/process-reservation/{src}/{id}/do", Repo.AdminProcessReservation)
	mux.Get("/admin/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservation)
	mux.Post("/admin/reservation/{src}/{id}", Repo.AdminPostShowReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

func CreateTestTemplate() (map[string]*template.Template, error) {
	cached := map[string]*template.Template{}
	//"./templates/*.page.html"

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplate))
	if err != nil {
		return cached, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("current file being processed is: ", name)

		ts, err := template.New(name).Funcs(funcmap).ParseFiles(page)

		if err != nil {
			return cached, err
		}

		//"./templates/*.layout.tmpl"

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))

		if err != nil {
			return cached, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return cached, err
			}
		}
		cached[name] = ts

	}

	return cached, nil
}

func WriteToConsole(hd http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//fmt.Println("page got hit!")
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


