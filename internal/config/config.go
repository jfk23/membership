package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/jfk23/gobookings/internal/model"
)

// AppConfig is for getting config
type AppConfig struct {
	UseCache       bool
	CachedTemplate map[string]*template.Template
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	InProduction   bool
	Session        *scs.SessionManager
	MailChan       chan model.MailData
}
