// Package application provides context to bind the uTeach app together.
package application

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // blank identifier import registers the sqlite driver
	"html/template"

	"github.com/umairidris/uTeach/config"
	"github.com/umairidris/uTeach/libtemplate"
)

// App is the application context which contains application-wide configuration and components.
type App struct {
	Config      *config.Config
	DB          *sqlx.DB
	CookieStore *sessions.CookieStore
	Templates   map[string]*template.Template
}

// New initializes a new App.
func New(configPath string) *App {
	conf := config.Load(configPath)

	db := sqlx.MustOpen("sqlite3", conf.DBPath)
	db.MustExec("PRAGMA foreign_keys=ON;")

	// cookie encryption key must be 32 bytes
	store := sessions.NewCookieStore([]byte(conf.CookieAuthenticationKey), []byte(conf.CookieEncryptionKey))
	templates := libtemplate.LoadTemplates(conf.TemplatesPath)
	app := &App{conf, db, store, templates}
	return app
}
