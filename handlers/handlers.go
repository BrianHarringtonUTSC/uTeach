// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/middleware"
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
	router := mux.NewRouter()
	m := middleware.Middleware{a}

	c := alice.New

	router.Handle("/", Handler{a, GetSubjects})
	router.Handle("/s/{subject}", Handler{a, GetThreads})
	router.Handle("/s/{subject}/submit", Handler{a, GetNewThread}).Methods("GET")
	router.Handle("/s/{subject}/submit", c(m.MustLogin).Then(Handler{a, PostNewThread})).Methods("POST")
	router.Handle("/user/{email}", Handler{a, GetUser})
	router.Handle("/login", Handler{a, GetLogin})
	router.Handle("/oauth2callback", Handler{a, GetOauth2Callback})
	router.Handle("/logout", Handler{a, GetLogout})

	t := c(m.MustLogin, m.SetThreadIDVar)
	router.Handle("/s/{subject}/{threadID}", c(m.SetThreadIDVar).Then(Handler{a, GetThread}))
	router.Handle("/t/{threadID}/upvote", t.Then(Handler{a, PostThreadVote})).Methods("POST")
	router.Handle("/t/{threadID}/upvote", t.Then(Handler{a, DeleteThreadVote})).Methods("DELETE")
	router.Handle("/t/{threadID}/hide", t.Append(m.MustBeAdminOrThreadCreator).Then(Handler{a, PostHideThread})).Methods("POST")
	router.Handle("/t/{threadID}/hide", t.Append(m.MustBeAdminOrThreadCreator).Then(Handler{a, DeleteHideThread})).Methods("DELETE")
	router.Handle("/t/{threadID}/pin", t.Append(m.MustBeAdmin).Then(Handler{a, PostPinThread})).Methods("POST")
	router.Handle("/t/{threadID}/pin", t.Append(m.MustBeAdmin).Then(Handler{a, DeletePinThread})).Methods("DELETE")

	// serve static files -- should be the last route
	staticFileServer := http.FileServer(http.Dir(a.Config.StaticFilesPath))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return router
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.App, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		// make sure user is nil so templates don't render a user
		data["SessionUser"] = nil
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

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
