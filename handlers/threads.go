package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/models"
)

// GetThreads renders all threads for the subject.
func GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := strings.ToLower(vars["subject"])

	app := application.GetFromContext(r)

	t := models.NewThreadModel(app.DB)

	data := map[string]interface{}{}

	// add pinned posts
	pinnedThreads, err := t.GetThreadsBySubjectAndIsPinned(subject, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	unpinnedThreads, err := t.GetThreadsBySubjectAndIsPinned(subject, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["PinnedThreads"] = pinnedThreads
	data["UnpinnedThreads"] = unpinnedThreads

	//  if there is a user, get the user's upvoted threads
	if user, ok := app.Store.SessionUser(r); ok {
		userUpvotedThreadIDs, err := t.GetThreadIdsUpvotedByEmail(user.Email)
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

	app := application.GetFromContext(r)
	t := models.NewThreadModel(app.DB)

	thread, err := t.GetThreadByID(threadID)
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
	subject := strings.ToLower(vars["subject"])

	app := application.GetFromContext(r)
	user, _ := app.Store.SessionUser(r)

	title := r.FormValue("title")
	text := r.FormValue("text")

	t := models.NewThreadModel(app.DB)
	thread, err := t.AddThread(title, text, subject, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, thread.URL(), 301)
}

// thread_vote is a helper for handling upvotes.
func thread_vote(w http.ResponseWriter, r *http.Request, upvoteFn func(int64, string) error) {
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

// PostThreadVote adds a vote for the useruser on a thread.
func PostThreadVote(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)
	t := models.NewThreadModel(app.DB)
	thread_vote(w, r, t.AddThreadVoteForUser)
}

// DeleteThreadVote removes a vote for the user on a thread.
func DeleteThreadVote(w http.ResponseWriter, r *http.Request) {
	app := application.GetFromContext(r)
	t := models.NewThreadModel(app.DB)
	thread_vote(w, r, t.RemoveTheadVoteForUser)
}

func PostHideThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.GetFromContext(r)

	t := models.NewThreadModel(app.DB)
	err = t.HideThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteHideThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.GetFromContext(r)

	t := models.NewThreadModel(app.DB)
	err = t.UnhideThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PostPinThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.GetFromContext(r)

	t := models.NewThreadModel(app.DB)
	err = t.PinThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePinThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := application.GetFromContext(r)

	t := models.NewThreadModel(app.DB)
	err = t.UnpinThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
