package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jfk23/gobookings/cmd/web/helpers"
	"github.com/jfk23/gobookings/driver"
	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/handlers"
	"github.com/jfk23/gobookings/internal/model"
	"github.com/jfk23/gobookings/internal/render"
)

const portNum = ":8080"
var appConfig config.AppConfig
var sessionManager *scs.SessionManager


func main() {
	db, err := run()
	if err !=nil {
		log.Fatal(err)
	}
	

	defer db.SQL.Close()

	defer close(appConfig.MailChan)

	ListenForMail()

	

	// send email using smtp standard mail package (not set-up localhost:1025 port yet!)
	// from := "hi@its.me"
	// auth := smtp.PlainAuth("", "me", "", "localhost")
	// err = smtp.SendMail("localhost:1025", auth, from, []string{"bye@its.you"}, []byte("Hi there!!"))
	// if err != nil {
	// 	log.Println(err)
	// }

	fmt.Println("server is running port on: ", portNum)
	// http.ListenAndServe(portNum, nil)

	srv := &http.Server{
		Addr: portNum,
		Handler: Routes(&appConfig),
	}

	e := srv.ListenAndServe()

	log.Fatal(e)

	

}

func run() (*driver.DB, error) {

	// below gob method is for storing Reservation data into session.
	gob.Register(model.Reservation{})
	gob.Register(model.User{})
	gob.Register(model.Room{})
	gob.Register(model.Restriction{})
	gob.Register(map[string]int{})

	inProduction :=flag.Bool("production", true, "default app run setting is production")
	useCache :=flag.Bool("cache", true, "default cache setting is use template cache")
	dbHost :=flag.String("dbhost", "localhost", "default db host is localhost")
	dbName :=flag.String("dbname", "", "default db name is bookings")
	dbUser :=flag.String("dbuser", "", "default db user is nuburi")
	dbPass :=flag.String("dbpass", "", "default db password is empty")
	dbPort :=flag.String("dbport", "5432", "default port number is 5432")
	dbSSL :=flag.String("dbssl", "disable", "ssl connection is either disable, prefer, require")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("missing required fields")
		os.Exit(1)
	}

	mailChan := make(chan model.MailData)

	appConfig.MailChan = mailChan

	appConfig.InProduction =*inProduction

	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = appConfig.InProduction
	sessionManager.Cookie.Persist = true

	appConfig.Session = sessionManager

	// connecting to database
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connect to database. bye....")
	} 
	log.Println("Now connected to database!")

	
	

	ts, err := render.CreateTemplate()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	appConfig.CachedTemplate = ts
	appConfig.UseCache = *useCache

	repo := handlers.CreateNewRepo(&appConfig, db)
	
	handlers.SetHandler(repo)

	render.NewRenderer(&appConfig)

	helpers.NewHelpers(&appConfig)

	return db, nil
}
