// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/middleware"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

// http://elithrar.github.io/article/http-handler-error-handling-revisited/
type Handler struct {
	App *application.App
	H   func(a *application.App, w http.ResponseWriter, r *http.Request) error
}

// Router gets the router with routes and their corresponding handlers defined.
// It also serves static files based on the static files path specified in the app config.
func Router(a *application.App) *mux.Router {

	// helper function to create Handler struct
	h := func(handlerFunc func(*application.App, http.ResponseWriter, *http.Request) error) *Handler {
		return &Handler{a, handlerFunc}
	}

	// app specific middleware
	m := middleware.Middleware{a}

	router := mux.NewRouter()
	router.Handle("/", h(GetSubjects))
	router.Handle("/s/{subject}", h(GetThreads))
	router.Handle("/s/{subject}/submit", h(GetNewThread)).Methods("GET")
	router.Handle("/s/{subject}/submit", m.MustLogin(h(PostNewThread))).Methods("POST")
	router.Handle("/user/{email}", h(GetUser))
	router.Handle("/login", h(GetLogin))
	router.Handle("/oauth2callback", h(GetOauth2Callback))
	router.Handle("/logout", h(GetLogout))

	// thread specific middleware
	t := alice.New(m.MustLogin, m.SetThreadIDVar)

	router.Handle("/t/{threadID}", m.SetThreadIDVar(h(GetThread)))
	router.Handle("/t/{threadID}/upvote", t.Then(h(PostThreadVote))).Methods("POST")
	router.Handle("/t/{threadID}/upvote", t.Then(h(DeleteThreadVote))).Methods("DELETE")
	router.Handle("/t/{threadID}/hide", t.Then(m.MustBeAdminOrThreadCreator(h(PostHideThread)))).Methods("POST")
	router.Handle("/t/{threadID}/hide", t.Then(m.MustBeAdminOrThreadCreator(h(DeleteHideThread)))).Methods("DELETE")
	router.Handle("/t/{threadID}/pin", t.Then(m.MustBeAdmin(h(PostPinThread)))).Methods("POST")
	router.Handle("/t/{threadID}/pin", t.Then(m.MustBeAdmin(h(DeletePinThread)))).Methods("DELETE")

	// serve static files -- should be the last route
	staticFileServer := http.FileServer(http.Dir(a.Config.StaticFilesPath))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return router
}

// ServeHTTP allows Handler to satisfy the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.App, w, r)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			http.NotFound(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println(err)
			debug.PrintStack()
		}
	}
}

// renderTemplate renders the template at name with data.
// It also adds the session user to the data for templates to access.
func renderTemplate(a *application.App, w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) error {
	tmpl, ok := a.Templates[name]
	if !ok {
		return errors.New(fmt.Sprintf("The template %s does not exist.", name))
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	// add session user to data
	usm := session.NewUserSessionManager(a.CookieStore)
	if user, ok := usm.SessionUser(r); ok {
		data["SessionUser"] = user
	} else {
		// pass in empty user
		data["SessionUser"] = &models.User{}
	}

	// TODO: to speed this up use a buffer pool (https://elithrar.github.io/article/using-buffer-pools-with-go/)
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}
	buf.WriteTo(w)
	return nil
}
