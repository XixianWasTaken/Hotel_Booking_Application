package render

import (
	"bytes"
	"github.com/justinas/nosurf"
	"html/template"
	"learningGo/cmd/internal/config"
	"learningGo/cmd/internal/modules"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

//NewTempaltes sets the config for the tempalte package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *modules.TemplateData, r *http.Request) *modules.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *modules.TemplateData, r *http.Request) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, checker := tc[tmpl]

	if !checker {
		log.Fatal("Could not get the template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}
	return myCache, err
}
