package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"learningGo/pkd/config"
	"learningGo/pkd/handlers"
	"learningGo/pkd/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

//main is the main application
func main() {
	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateTest()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	serving := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serving.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprint("Starting on application on port: ", portNumber))

}
