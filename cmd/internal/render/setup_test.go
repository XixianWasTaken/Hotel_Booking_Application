package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"learningGo/cmd/internal/config"
	"learningGo/cmd/internal/modules"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testapp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(modules.Reservation{})

	//change this to true when in production
	testapp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testapp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testapp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testapp.Session = session

	app = &testapp
	os.Exit(m.Run())
}

type myWriter struct{}

func (w *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (w *myWriter) WriteHeader(i int) {

}

func (w *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
