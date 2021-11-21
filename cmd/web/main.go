package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"learningGo/cmd/internal/config"
	"learningGo/cmd/internal/driver"
	"learningGo/cmd/internal/handlers"
	"learningGo/cmd/internal/helpers"
	"learningGo/cmd/internal/modules"
	"learningGo/cmd/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//main is the main application
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	log.Println(fmt.Sprint("Starting on application on port: ", portNumber))

	serving := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serving.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	//what am I going to put in the session
	gob.Register(modules.Reservation{})
	gob.Register(modules.User{})
	gob.Register(modules.Room{})
	gob.Register(modules.Restrictions{})

	//change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database....")
	db, err := driver.ConnectSQL("host = localhost port = 5432 dbname = bookings user = postgres password = 123456")
	if err != nil {
		log.Fatal("Cannot connect to database!")
	}

	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
