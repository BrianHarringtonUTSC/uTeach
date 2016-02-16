// Package routes provides route handlers for the uTeach app.
package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/umairidris/uTeach/app"
)

type RouteHandler struct {
	App *app.App
}

func (rh *RouteHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := rh.App.DB.Subjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	rh.App.RenderTemplate(w, r, "subjects.html", data)
}

func (rh *RouteHandler) GetTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]

	topics, err := rh.App.DB.Topics(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Topics": topics}
	rh.App.RenderTemplate(w, r, "topics.html", data)
}

func (rh *RouteHandler) GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]

	threads, err := rh.App.DB.Threads(subjectName, topicName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Threads": threads}

	if user, ok := rh.App.Store.SessionUser(r); ok {
		userUpvotedThreadIDs, err := rh.App.DB.UserUpvotedThreadIDs(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["UserUpvotedThreadIDs"] = userUpvotedThreadIDs
	}

	rh.App.RenderTemplate(w, r, "threads.html", data)
}

func (rh *RouteHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := rh.App.DB.Thread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Thread": thread}
	rh.App.RenderTemplate(w, r, "thread.html", data)
}

func (rh *RouteHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	userCreatedThreads, err := rh.App.DB.UserCreatedThreads(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Username": username, "UserCreatedThreads": userCreatedThreads}
	rh.App.RenderTemplate(w, r, "user.html", data)
}

func (rh *RouteHandler) Login(w http.ResponseWriter, r *http.Request) {
	if _, ok := rh.App.Store.SessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	// TODO: replace this with SAML login for username

	err := rh.App.Store.NewUserSession(w, r, username, rh.App.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged in as: "+username)
}

func (rh *RouteHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := rh.App.Store.DeleteUserSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged out")
}

func (rh *RouteHandler) upvote(w http.ResponseWriter, r *http.Request, upvoteFn func(string, int) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := rh.App.Store.SessionUser(r)
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

func (rh *RouteHandler) AddUpvote(w http.ResponseWriter, r *http.Request) {
	rh.upvote(w, r, rh.App.DB.AddUpVote)
}

func (rh *RouteHandler) RemoveUpvote(w http.ResponseWriter, r *http.Request) {
	rh.upvote(w, r, rh.App.DB.RemoveUpvote)
}
