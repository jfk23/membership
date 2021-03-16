package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jfk23/gobookings/pkg/config"
	"github.com/jfk23/gobookings/pkg/handlers"
	"github.com/jfk23/gobookings/pkg/render"
)

const portNum = ":8080"
var appConfig config.AppConfig
var sessionManager *scs.SessionManager


func main() {

	appConfig.InProduction =false

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = appConfig.InProduction
	sessionManager.Cookie.Persist = true

	appConfig.Session = sessionManager
	

	ts, err := render.CreateTemplate()
	if err != nil {
		log.Fatal(err)
	}
	appConfig.CachedTemplate = ts
	appConfig.UseCache = false

	repo := handlers.CreateNewRepo(&appConfig)
	
	handlers.SetHandler(repo)

	render.LoadTemplate(&appConfig)

	// fmt.Println(ts)

	// fmt.Printf("type of ts is %T\n", ts)
	// fmt.Printf("type of t.cachedtemplate is %T\n", T.CachedTemplate)



	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("server is running port on: ", portNum)
	// http.ListenAndServe(portNum, nil)

	srv := &http.Server{
		Addr: portNum,
		Handler: Routes(&appConfig),
	}

	e := srv.ListenAndServe()

	log.Fatal(e)

}
