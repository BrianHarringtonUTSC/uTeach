// Package handlers provides route handlers for the uTeach app.
package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
	"strconv"

	"github.com/UmairIdris/uTeach/application"
	"github.com/UmairIdris/uTeach/middleware"
)

// Router gets the router with routes and their corresponding handlers defined.
// It also serves static files based on the static files path specified in the app config.
func Router(app *application.Application) *mux.Router {
	stdChain := alice.New(middleware.SetApplication(app))
	authChain := stdChain.Append(middleware.MustLogin)

	router := mux.NewRouter()
	router.Handle("/", stdChain.ThenFunc(GetSubjects))
	router.Handle("/threads/{subjectName}", stdChain.ThenFunc(GetThreads))
	router.Handle("/thread/{subjectName}/{threadID}", stdChain.ThenFunc(GetThread))
	outer.Handle("/thread/{subjectName}/submit", stdChain.ThenFunc(GetNewThread))

	router.Handle("/user/{username}", stdChain.ThenFunc(GetUser))

	router.Handle("/login/{username}", stdChain.ThenFunc(Login))
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
	subjectName := vars["subjectName"]

	app := application.Get(r)
	threads, err := app.DB.Threads(subjectName)
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
	threadID, err := strconv.Atoi(vars["threadID"])
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

func GetNewThread(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "new_thread.html", nil)
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

// Login logs the user in.
func Login(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)

	if _, ok := app.Store.SessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	// TODO: replace this with SAML login for username

	err := app.Store.NewUserSession(w, r, username, app.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged in as: "+username)
}

// Logout logs the user out.
func Logout(w http.ResponseWriter, r *http.Request) {
	app := application.Get(r)
	err := app.Store.DeleteUserSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged out")
}

// upvote is a helper for handling upvotes.
func upvote(w http.ResponseWriter, r *http.Request, upvoteFn func(string, int) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.Get(r)
	user, ok := app.Store.SessionUser(r)
	if !ok {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	err = upvoteFn(user.Username, threadID)
	if err != nil {
		fmt.Println(err)
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
