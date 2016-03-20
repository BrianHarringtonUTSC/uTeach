// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	m "github.com/umairidris/uTeach/middleware"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

// Router gets the router with routes and their corresponding handlers defined.
// It also serves static files based on the static files path specified in the app config.
func Router(app *application.App) *mux.Router {
	router := mux.NewRouter()

	c := alice.New(m.SetApplication(app))
	router.Handle("/", c.ThenFunc(GetSubjects))
	router.Handle("/s/{subject}", c.ThenFunc(GetThreads))
	router.Handle("/s/{subject}/submit", c.ThenFunc(GetNewThread)).Methods("GET")
	router.Handle("/s/{subject}/submit", c.Append(m.MustLogin).ThenFunc(PostNewThread)).Methods("POST")
	router.Handle("/user/{email}", c.ThenFunc(GetUser))
	router.Handle("/login", c.ThenFunc(GetLogin))
	router.Handle("/oauth2callback", c.ThenFunc(GetOauth2Callback))
	router.Handle("/logout", c.ThenFunc(Logout))

	t := c.Append(m.MustLogin, m.SetThreadIDVar)
	router.Handle("/s/{subject}/{threadID}", c.Append(m.SetThreadIDVar).ThenFunc(GetThread))
	router.Handle("/t/{threadID}/upvote", t.ThenFunc(PostThreadVote)).Methods("POST")
	router.Handle("/t/{threadID}/upvote", t.ThenFunc(DeleteThreadVote)).Methods("DELETE")
	router.Handle("/t/{threadID}/hide", t.Append(m.MustBeAdminOrThreadCreator).ThenFunc(PostHideThread)).Methods("POST")
	router.Handle("/t/{threadID}/hide", t.Append(m.MustBeAdminOrThreadCreator).ThenFunc(DeleteHideThread)).Methods("DELETE")
	router.Handle("/t/{threadID}/pin", t.Append(m.MustBeAdmin).ThenFunc(PostPinThread)).Methods("POST")
	router.Handle("/t/{threadID}/pin", t.Append(m.MustBeAdmin).ThenFunc(DeletePinThread)).Methods("DELETE")

	// serve static files -- should be the last route
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(app.Config.StaticFilesPath))))

	return router
}

// renderTemplate renders the template at name with data.
// It also adds the session user to the data for templates to access.
func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	templates := context.Templates(r)
	tmpl, ok := templates[name]
	if !ok {
		handleError(w, errors.New(fmt.Sprintf("The template %s does not exist.", name)))
		return
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	// add session user to data
	if user, ok := getSessionUser(r); ok {
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

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getSessionUser(r *http.Request) (*models.User, bool) {
	usm := session.NewUserSessionManager(context.CookieStore(r))
	return usm.SessionUser(r)
}
