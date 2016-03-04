// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/middleware"
)

// Router gets the router with routes and their corresponding handlers defined.
// It also serves static files based on the static files path specified in the app config.
func Router(app *application.Application) *mux.Router {
	stdChain := alice.New(middleware.SetApplication(app))
	authChain := stdChain.Append(middleware.MustLogin)
	// adminChain := authChain.Append(middleware.MustBeAdmin)

	router := mux.NewRouter()
	router.Handle("/", stdChain.ThenFunc(GetSubjects))
	router.Handle("/s/{subject}", stdChain.ThenFunc(GetThreads))
	router.Handle("/s/{subject}/submit", authChain.ThenFunc(GetNewThread)).Methods("GET")
	router.Handle("/s/{subject}/submit", authChain.ThenFunc(PostNewThread)).Methods("POST")
	router.Handle("/s/{subject}/{threadID}", stdChain.ThenFunc(GetThread))

	router.Handle("/t/{threadID}/upvote", authChain.ThenFunc(PostThreadVote)).Methods("POST")
	router.Handle("/t/{threadID}/upvote", authChain.ThenFunc(DeleteThreadVote)).Methods("DELETE")
	router.Handle("/t/{threadID}/hide", authChain.ThenFunc(PostHideThread)).Methods("POST")
	router.Handle("/t/{threadID}/hide", authChain.ThenFunc(DeleteHideThread)).Methods("DELETE")
	router.Handle("/t/{threadID}/pin", authChain.ThenFunc(PostPinThread)).Methods("POST")
	router.Handle("/t/{threadID}/pin", authChain.ThenFunc(DeletePinThread)).Methods("DELETE")

	router.Handle("/user/{email}", stdChain.ThenFunc(GetUser))

	router.Handle("/login", stdChain.ThenFunc(GetLogin))
	router.Handle("/oauth2callback", stdChain.ThenFunc(GetOauth2Callback))
	router.Handle("/logout", stdChain.ThenFunc(Logout))

	// serve static files -- should be the last route
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(app.Config.StaticFilesPath))))

	return router
}

// renderTemplate renders the template at name with data.
// It also adds the session user to the data for templates to access.
func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	app := application.GetFromContext(r)

	tmpl, ok := app.Templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name), http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	// add session user to data
	if user, ok := app.Store.SessionUser(r); ok {
		data["SessionUser"] = user
		data["IsAdmin"] = user.IsAdmin
	} else {
		// make sure user is nil so templates don't render a user
		data["SessionUser"] = nil
		data["IsAdmin"] = false
	}

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}
