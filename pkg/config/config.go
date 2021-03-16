package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig is for getting config
type AppConfig struct {
	UseCache       bool
	CachedTemplate map[string]*template.Template
	InfoLog        *log.Logger
	InProduction   bool
	Session        *scs.SessionManager
}
