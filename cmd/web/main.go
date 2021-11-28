package main

import (
	"encoding/gob"
	"flag"
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
	defer close(app.MailChan)

	listenForMail()

	log.Println(fmt.Sprint("Starting on application on port ", portNumber))
	log.Println("Starting mail listener...")

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
	gob.Register(map[string]int{})

	/*read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost","localhost","Database host")
	dbName := flag.String("dbname","","Database name")
	dbUser := flag.String("dbuser","","Database user")
	dbPass := flag.String("dbpass","","Database password")
	dbPort := flag.String("dbport","5432","Database port")
	dbSSL := flag.String("dbssl","disable","Database ssl setting (disable, prefer, require)")

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}
	*/

	flag.Parse()

	mailChan := make(chan modules.MailData)
	app.MailChan = mailChan

	/*
		app.InProduction = *inProduction
		app.UseCache = *useCache
	*/

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
	/*
		connectionString := fmt.Sprintf("host = %s port = %s dbname = %s user = %s password = %s sslmode = %s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	*/
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
