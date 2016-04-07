// Package application provides context to bind the app together.
package application

import (
	"html/template"
	"log"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // blank identifier import registers the sqlite driver
	"github.com/umairidris/uTeach/config"
	"github.com/umairidris/uTeach/libtemplate"
)

// App is the context which contains application-wide configuration and components.
type App struct {
	Config    *config.Config
	DB        *sqlx.DB
	Store     sessions.Store
	Templates map[string]*template.Template
}

// New initializes a new App.
func New(configPath string) *App {
	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	db := sqlx.MustOpen("sqlite3", conf.DBPath)
	db.MustExec("PRAGMA foreign_keys=ON;")

	// cookie encryption key must be 32 bytes
	store := sessions.NewCookieStore([]byte(conf.CookieAuthenticationKey), []byte(conf.CookieEncryptionKey))

	templates, err := libtemplate.Load(conf.TemplatesPath)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{conf, db, store, templates}
	return app
}
