// Package application provides context to bind the app together.
package application

import (
	"encoding/base64"
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

// New creates a new App based on the config. Exits if an error is encountered.
func New(conf config.Config) *App {
	db := sqlx.MustOpen("sqlite3", conf.DBPath)
	db.MustExec("PRAGMA foreign_keys=ON;")

	cookieAuthKey, err := base64.StdEncoding.DecodeString(conf.CookieAuthenticationKeyBase64)
	if err != nil {
		log.Fatal(err)
	}

	cookieEncKey, err := base64.StdEncoding.DecodeString(conf.CookieEncryptionKeyBase64)
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore(cookieAuthKey, cookieEncKey)

	templates, err := libtemplate.Load(conf.TemplatesPath)
	if err != nil {
		log.Fatal(err)
	}

	return &App{&conf, db, store, templates}
}
