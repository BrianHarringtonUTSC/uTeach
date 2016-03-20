package context

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"html/template"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/config"
)

const (
	appKey = "app"

	// app related keys
	configKey      = "config"
	cookieStoreKey = "cookiestore"
	dbKey          = "db"
	templatesKey   = "templates"

	// handler related keys
	threadIDKey = "threadID"
)

// SetContext sets the global application in the request.
func SetApp(r *http.Request, a *application.App) {
	context.Set(r, appKey, a)
	context.Set(r, configKey, a.Config)
	context.Set(r, cookieStoreKey, a.CookieStore)
	context.Set(r, dbKey, a.DB)
	context.Set(r, templatesKey, a.Templates)
}

func Config(r *http.Request) *config.Config {
	return context.Get(r, configKey).(*config.Config)
}

func CookieStore(r *http.Request) *sessions.CookieStore {
	return context.Get(r, cookieStoreKey).(*sessions.CookieStore)
}

func DB(r *http.Request) *sqlx.DB {
	return context.Get(r, dbKey).(*sqlx.DB)
}

func Templates(r *http.Request) map[string]*template.Template {
	return context.Get(r, templatesKey).(map[string]*template.Template)
}

func SetThreadID(r *http.Request, threadID int64) {
	context.Set(r, threadIDKey, threadID)
}

func ThreadID(r *http.Request) int64 {
	return context.Get(r, threadIDKey).(int64)
}
