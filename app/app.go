// Package app provides context to bind the uTeach app together.
package app

import (
	"html/template"

	"github.com/umairidris/uTeach/config"
	"github.com/umairidris/uTeach/db"
	"github.com/umairidris/uTeach/session"
)

// App is the application context, the central struct binding all components of the app together.
type App struct {
	Config    *config.Config
	DB        *db.DB
	Store     *session.Store
	Templates map[string]*template.Template
}

// New initializes a new App.
func New(configPath string) *App {
	config := config.Load(configPath)
	db := db.New(config.DBPath)
	store := session.NewStore("todo-proper-secret")
	templates := LoadTemplates(config.TemplatesPath)

	app := &App{config, db, store, templates}

	return app
}
