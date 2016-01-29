package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleGetSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := GetSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Subjects": subjects}
	RenderTemplate(w, r, "subjects.html", data)
}

func handleGetTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]

	topics, err := GetTopics(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Topics": topics}
	RenderTemplate(w, r, "topics.html", data)
}

func handleGetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]

	threads, err := GetThreads(subjectName, topicName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userUpvotedThreadIDs := make(map[int]bool)
	user, ok := GetSessionUser(r)
	if ok {
		userUpvotedThreadIDs, err = GetUserUpvotedThreadIDs(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data := map[string]interface{}{"Threads": threads, "UserUpvotedThreadIDs": userUpvotedThreadIDs}
	RenderTemplate(w, r, "threads.html", data)
}

func handleGetThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := GetThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Thread": thread}
	RenderTemplate(w, r, "thread.html", data)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := GetSessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	// TODO: replace this with SAML login for username

	err := NewUserSession(w, r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged in as: "+username)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	err := DeleteUserSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Logged out")
}

func handleUpvote(w http.ResponseWriter, r *http.Request, fn func(string, int) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := GetSessionUser(r)
	if !ok {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	err = fn(user.Username, threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func handleAddUpvote(w http.ResponseWriter, r *http.Request) {
	handleUpvote(w, r, AddUpVote)
}

func handleRemoveUpvote(w http.ResponseWriter, r *http.Request) {
	handleUpvote(w, r, RemoveUpvote)
}
