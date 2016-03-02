// Package application provides context to bind the uTeach app together.
package application

import (
	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // blank identifier import registers the sqlite driver
	"html/template"
	"net/http"

	"github.com/umairidris/uTeach/config"
	"github.com/umairidris/uTeach/libtemplate"
	"github.com/umairidris/uTeach/session"
)

const (
	contextKey = "application" // key for storing/retrieving Application from gorilla context
)

// Application is the application context which contains application-wide configuration and components.
type Application struct {
	Config    *config.Config
	DB        *sqlx.DB
	Store     *session.Store
	Templates map[string]*template.Template
}

// New initializes a new Application.
func New(configPath string) *Application {
	config := config.Load(configPath)

	db := sqlx.MustOpen("sqlite3", config.DBPath)
	db.MustExec("PRAGMA foreign_keys=ON;")

	// cookie encryption key must be 32 bytes
	store := session.NewStore(config.CookieAuthenticationKey, config.CookieEncryptionKey)
	templates := libtemplate.LoadTemplates(config.TemplatesPath)
	app := &Application{config, db, store, templates}
	return app
}

// SetContext sets the global application in the request.
func SetInContext(r *http.Request, app *Application) {
	context.Set(r, contextKey, app)
}

// Get gets the global application from the request. Application MUST be set in the context before, else the application
// will be nil. See SetContext helper function.
func GetFromContext(r *http.Request) *Application {
	return context.Get(r, contextKey).(*Application)
}
