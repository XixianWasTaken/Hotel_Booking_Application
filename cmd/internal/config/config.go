package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"learningGo/cmd/internal/modules"
	"log"
)

//AppConfig holds hte application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan modules.MailData
}
