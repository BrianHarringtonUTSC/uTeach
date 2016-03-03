package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

// upvote is a helper for handling upvotes.
func upvote(w http.ResponseWriter, r *http.Request, upvoteFn func(int64, string) error) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.GetFromContext(r)
	user, _ := app.Store.SessionUser(r)

	err = upvoteFn(threadID, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddUpvote adds an upvote for the useruser on a thread.
func AddUpvote(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)
	t := models.NewThreadModel(app.DB)
	upvote(w, r, t.AddThreadVoteForUser)
}

// RemoveUpvote removes an upvote for the user on a thread.
func RemoveUpvote(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)
	t := models.NewThreadModel(app.DB)
	upvote(w, r, t.RemoveTheadVoteForUser)
}
