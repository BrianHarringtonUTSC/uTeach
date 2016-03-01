// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
	"strconv"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/middleware"
)

// Router gets the router with routes and their corresponding handlers defined.
// It also serves static files based on the static files path specified in the app config.
func Router(app *application.Application) *mux.Router {
	stdChain := alice.New(middleware.SetApplication(app))
	authChain := stdChain.Append(middleware.MustLogin)

	router := mux.NewRouter()
	router.Handle("/", stdChain.ThenFunc(GetSubjects))
	router.Handle("/s/{subject}", stdChain.ThenFunc(GetThreads))
	router.Handle("/s/{subject}/submit", authChain.ThenFunc(GetNewThread)).Methods("GET")
	router.Handle("/s/{subject}/submit", authChain.ThenFunc(PostNewThread)).Methods("POST")
	router.Handle("/s/{subject}/{threadID}", stdChain.ThenFunc(GetThread))

	router.Handle("/user/{username}", stdChain.ThenFunc(GetUser))

	router.Handle("/login", stdChain.ThenFunc(GetLogin))
	router.Handle("/oauth2callback", stdChain.ThenFunc(GetOauth2Callback))

	router.Handle("/logout", stdChain.ThenFunc(Logout))

	router.Handle("/upvote/{threadID}", authChain.ThenFunc(AddUpvote)).Methods("POST")
	router.Handle("/upvote/{threadID}", authChain.ThenFunc(RemoveUpvote)).Methods("DELETE")

	// serve static files -- should be the last route
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(app.Config.StaticFilesPath))))

	return router
}

// renderTemplate renders the template at name with data.
// It also adds the session user to the data for templates to access.
func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	app := application.Get(r)

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
	} else {
		// make sure user is nil so templates don't render a user
		data["SessionUser"] = nil
	}

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}

// GetSubjects renders all subjects.
func GetSubjects(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	subjects, err := app.DB.Subjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	renderTemplate(w, r, "subjects.html", data)
}

// GetThreads renders all threads for the subject.
func GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"]

	app := application.Get(r)
	threads, err := app.DB.Threads(subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Threads": threads}

	if user, ok := app.Store.SessionUser(r); ok {
		userUpvotedThreadIDs, err := app.DB.UserUpvotedThreadIDs(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["UserUpvotedThreadIDs"] = userUpvotedThreadIDs
	}

	renderTemplate(w, r, "threads.html", data)
}

// GetThread renders a thread.
func GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.Get(r)
	thread, err := app.DB.Thread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Thread": thread}
	renderTemplate(w, r, "thread.html", data)
}

// GetNewThread renders the new thread page.
func GetNewThread(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "new_thread.html", nil)
}

// PostNewThread adds a new thread in the db and redirects to it, if successful.
func PostNewThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"]

	app := application.Get(r)
	user, _ := app.Store.SessionUser(r)

	title := r.FormValue("title")
	text := r.FormValue("text")

	thread, err := app.DB.NewThread(title, text, subject, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, thread.URL(), 301)
}

// GetUser renders user info.
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	app := application.Get(r)
	userCreatedThreads, err := app.DB.UserCreatedThreads(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Username": username, "UserCreatedThreads": userCreatedThreads}
	renderTemplate(w, r, "user.html", data)
}

// upvote is a helper for handling upvotes.
func upvote(w http.ResponseWriter, r *http.Request, upvoteFn func(string, int64) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.Get(r)
	user, _ := app.Store.SessionUser(r)

	err = upvoteFn(user.Username, threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddUpvote adds an upvote for the user on a thread.
func AddUpvote(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	upvote(w, r, app.DB.AddUpVote)
}

// RemoveUpvote removes an upvote for the user on a thread.
func RemoveUpvote(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	upvote(w, r, app.DB.RemoveUpvote)
}
