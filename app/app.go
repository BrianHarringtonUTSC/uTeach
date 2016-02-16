// Package app provides context to bind the uTeach app together.
package app

import (
	"html/template"

	"github.com/umairidris/uTeach/db"
	"github.com/umairidris/uTeach/session"
)

// App is the application context, the central struct binding all components of the app together.
type App struct {
	DB        *db.DB
	Store     *session.Store
	Templates map[string]*template.Template
}

// NewApp initializes a new App.
func New() *App {
	db := db.New("./uteach.db")
	store := session.NewStore("todo-proper-secret")
	templates := LoadTemplates("tmpl/")

	app := &App{db, store, templates}

	return app
}
