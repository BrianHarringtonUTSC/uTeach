// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/httperror"
	"github.com/umairidris/uTeach/libtemplate"
	"github.com/umairidris/uTeach/middleware"
	"github.com/umairidris/uTeach/models"
)

// Handler struct used to pass application context to a handler and get an error for cleaner error handling.
// See: http://elithrar.github.io/article/http-handler-error-handling-revisited/
type Handler struct {
	App *application.App
	H   func(a *application.App, w http.ResponseWriter, r *http.Request) error
}

// Router gets the router with routes and their corresponding handlers defined.
// It also serves static files based on the static files path specified in the app config.
func Router(a *application.App) http.Handler {
	// helper function to create Handler struct
	h := func(handlerFunc func(*application.App, http.ResponseWriter, *http.Request) error) http.Handler {
		return &Handler{a, handlerFunc}
	}

	// app specific middleware
	m := middleware.Middleware{a}

	// middleware for all routes
	standardChain := alice.New(m.SetSessionUser)

	router := mux.NewRouter()

	// subject routes
	router.Handle("/", h(getSubjects))

	// user routes
	router.Handle("/user/{email}", h(getUser))
	router.Handle("/login", h(getLogin))
	router.Handle("/oauth2callback", h(getOauth2Callback))
	router.Handle("/logout", h(getLogout))

	// thread routes
	t := alice.New(m.MustLogin, m.SetThread)
	router.Handle("/s/{subject}", m.SetSubject(h(getThreads)))
	router.Handle("/s/{subject}/new", m.SetSubject(h(getNewThread))).Methods("GET")
	router.Handle("/s/{subject}/new", m.MustLogin(m.SetSubject(h(postNewThread)))).Methods("POST")
	router.Handle("/t/{threadID}", m.SetThread(h(getThread)))
	router.Handle("/t/{threadID}/upvote", t.Then(h(postThreadVote))).Methods("POST")
	router.Handle("/t/{threadID}/upvote", t.Then(h(deleteThreadVote))).Methods("DELETE")
	router.Handle("/t/{threadID}/hide", t.Then(m.MustBeAdminOrThreadCreator(h(postHideThread)))).Methods("POST")
	router.Handle("/t/{threadID}/hide", t.Then(m.MustBeAdminOrThreadCreator(h(deleteHideThread)))).Methods("DELETE")
	router.Handle("/t/{threadID}/pin", t.Then(m.MustBeAdmin(h(postPinThread)))).Methods("POST")
	router.Handle("/t/{threadID}/pin", t.Then(m.MustBeAdmin(h(deletePinThread)))).Methods("DELETE")

	// tag routes
	router.Handle("/s/{subject}/tags", m.SetSubject(h(getTags)))
	router.Handle("/s/{subject}/tags/{tag}", m.SetSubject(m.SetTag(h(getThreadsByTag))))

	// serve static files -- should be the last route
	staticFileServer := http.FileServer(http.Dir(a.Config.StaticFilesPath))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return standardChain.Then(router)
}

// ServeHTTP allows Handler to satisfy the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.App, w, r)
	httperror.HandleError(w, err)
}

// renderTemplate renders the template at name with data.
// It also adds the session user to the data for templates to access.
func renderTemplate(a *application.App, w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) error {
	// add session user to data
	if user, ok := context.SessionUser(r); ok {
		data["SessionUser"] = user
	} else {
		// pass in empty user
		data["SessionUser"] = &models.User{}
	}

	return libtemplate.Render(w, a.Templates, name, data)
}
