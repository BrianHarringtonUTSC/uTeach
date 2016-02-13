package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (a *App) handleGetSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := a.db.Subjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	a.RenderTemplate(w, r, "subjects.html", data)
}

func (a *App) handleGetTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]

	topics, err := a.db.Topics(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Topics": topics}
	a.RenderTemplate(w, r, "topics.html", data)
}

func (a *App) handleGetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]

	threads, err := a.db.Threads(subjectName, topicName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userUpvotedThreadIDs := make(map[int]bool)
	if user, ok := a.store.SessionUser(r); ok {
		userUpvotedThreadIDs, err = a.db.UserUpvotedThreadIDs(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data := map[string]interface{}{"Threads": threads, "UserUpvotedThreadIDs": userUpvotedThreadIDs}
	a.RenderTemplate(w, r, "threads.html", data)
}

func (a *App) handleGetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := a.db.Thread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Thread": thread}
	a.RenderTemplate(w, r, "thread.html", data)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.store.SessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	// TODO: replace this with SAML login for username

	err := a.store.NewUserSession(w, r, username, a.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged in as: "+username)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	err := a.store.DeleteUserSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged out")
}

func (a *App) handleUpvote(w http.ResponseWriter, r *http.Request, upvoteFn func(string, int) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := a.store.SessionUser(r)
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

func (a *App) handleAddUpvote(w http.ResponseWriter, r *http.Request) {
	a.handleUpvote(w, r, a.db.AddUpVote)
}

func (a *App) handleRemoveUpvote(w http.ResponseWriter, r *http.Request) {
	a.handleUpvote(w, r, a.db.RemoveUpvote)
}
