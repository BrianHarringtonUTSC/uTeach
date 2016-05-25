// Package handlers provides a router with routes attached.
package handlers

import (
	"net/http"

	"github.com/BrianHarringtonUTSC/uTeach/application"
	"github.com/BrianHarringtonUTSC/uTeach/httperror"
	"github.com/BrianHarringtonUTSC/uTeach/middleware"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// Handler allows passing an application context to a handler and handling errors.
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

	router := mux.NewRouter()
	router.StrictSlash(true)

	// topic routes
	router.Handle("/", h(getTopics))
	router.Handle("/topics/new", m.MustBeAdmin(h(getNewTopic))).Methods("GET")
	router.Handle("/topics/new", m.MustBeAdmin(h(postNewTopic))).Methods("POST")

	// user routes
	router.Handle("/users/{email}", h(getUser))
	router.Handle("/login", h(getLogin))
	router.Handle("/oauth2callback", h(getOauth2Callback))
	router.Handle("/logout", h(getLogout))

	// tag routes
	router.Handle("/topics/{topicName}/tags", m.SetTopic(h(getTags)))
	router.Handle("/topics/{topicName}/tags/new", m.MustBeAdmin(m.SetTopic(h(getNewTag)))).Methods("GET")
	router.Handle("/topics/{topicName}/tags/new", m.MustBeAdmin(m.SetTopic(h(postNewTag)))).Methods("POST")
	router.Handle("/topics/{topicName}/tags/{tagName}", m.SetTopic(m.SetTag(h(getPostsByTag))))

	// post routes
	p := alice.New(m.SetTopic)
	router.Handle("/topics/{topicName}", p.Then(h(getPosts)))

	p = p.Append(m.MustLogin)
	router.Handle("/topics/{topicName}/new", p.Then(h(getNewPost))).Methods("GET")
	router.Handle("/topics/{topicName}/new", p.Then(m.SetTopic(h(postNewPost)))).Methods("POST")

	p = p.Append(m.SetPost)
	router.Handle("/topics/{topicName}/posts/{postID}", m.SetTopic(m.SetPost(h(getPost))))
	router.Handle("/topics/{topicName}/posts/{postID}/vote", p.Then(h(postPostVote))).Methods("POST")
	router.Handle("/topics/{topicName}/posts/{postID}/vote", p.Then(h(deletePostVote))).Methods("DELETE")
	router.Handle("/topics/{topicName}/posts/{postID}/hide", p.Then(m.MustBeAdminOrPostCreator(h(postHidePost)))).Methods("POST")
	router.Handle("/topics/{topicName}/posts/{postID}/hide", p.Then(m.MustBeAdminOrPostCreator(h(deleteHidePost)))).Methods("DELETE")
	router.Handle("/topics/{topicName}/posts/{postID}/pin", p.Then(m.MustBeAdmin(h(postPinPost)))).Methods("POST")
	router.Handle("/topics/{topicName}/posts/{postID}/pin", p.Then(m.MustBeAdmin(h(deletePinPost)))).Methods("DELETE")

	// serve static files -- should be the last route
	staticFileServer := http.FileServer(http.Dir(a.Config.StaticFilesPath))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	// middleware for all routes
	standardChain := alice.New(m.SetTemplateData, m.SetSessionUser)
	return standardChain.Then(router)
}

// ServeHTTP allows Handler to satisfy the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.App, w, r)
	httperror.HandleError(w, err)
}
