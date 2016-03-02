package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/umairidris/uTeach/application"
)

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
