package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"

	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

func getThreadModel(r *http.Request) *models.ThreadModel {
	return models.NewThreadModel(context.DB(r))
}

// GetThreads renders all threads for the subject.
func GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := strings.ToLower(vars["subject"])

	// TODO: check if subject exists
	tm := getThreadModel(r)

	pinnedThreads, err := tm.GetThreadsBySubjectAndIsPinned(subject, true)
	if err != nil {
		handleError(w, err)
		return
	}

	unpinnedThreads, err := tm.GetThreadsBySubjectAndIsPinned(subject, false)
	if err != nil {
		handleError(w, err)
		return
	}

	data := map[string]interface{}{}
	data["PinnedThreads"] = pinnedThreads
	data["UnpinnedThreads"] = unpinnedThreads

	//  if there is a user, get the user's upvoted threads
	if user, ok := getSessionUser(r); ok {
		userUpvotedThreadIDs, err := tm.GetThreadIdsUpvotedByEmail(user.Email)
		if err != nil {
			handleError(w, err)
			return
		}
		data["UserUpvotedThreadIDs"] = userUpvotedThreadIDs
	}

	renderTemplate(w, r, "threads.html", data)
}

// GetThread renders a thread.
func GetThread(w http.ResponseWriter, r *http.Request) {
	tm := models.NewThreadModel(context.DB(r))

	threadID := context.ThreadID(r)
	thread, err := tm.GetThreadByID(threadID)
	if err != nil {
		handleError(w, err)
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

	user, _ := getSessionUser(r)

	title := r.FormValue("title")
	text := r.FormValue("text")

	tm := models.NewThreadModel(context.DB(r))
	thread, err := tm.AddThread(title, text, subject, user.Email)
	if err != nil {
		handleError(w, err)
		return
	}
	http.Redirect(w, r, thread.URL(), 301)
}

func handleThreadAction(w http.ResponseWriter, r *http.Request, f func(int64) error) {
	threadID := context.ThreadID(r)

	if err := f(threadID); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteThreadVote removes a vote for the user on a thread.
func PostThreadVote(w http.ResponseWriter, r *http.Request) {
	user, _ := getSessionUser(r)
	tm := getThreadModel(r)

	f := func(id int64) error {
		return tm.AddThreadVoteForUser(id, user.Email)
	}

	handleThreadAction(w, r, f)
}

// DeleteThreadVote removes a vote for the user on a thread.
func DeleteThreadVote(w http.ResponseWriter, r *http.Request) {
	user, _ := getSessionUser(r)
	tm := getThreadModel(r)

	f := func(id int64) error {
		return tm.RemoveTheadVoteForUser(id, user.Email)
	}

	handleThreadAction(w, r, f)
}

func PostHideThread(w http.ResponseWriter, r *http.Request) {
	tm := getThreadModel(r)
	handleThreadAction(w, r, tm.HideThread)
}

func DeleteHideThread(w http.ResponseWriter, r *http.Request) {
	tm := getThreadModel(r)
	handleThreadAction(w, r, tm.UnhideThread)
}

func PostPinThread(w http.ResponseWriter, r *http.Request) {
	tm := getThreadModel(r)
	handleThreadAction(w, r, tm.PinThread)
}

func DeletePinThread(w http.ResponseWriter, r *http.Request) {
	tm := getThreadModel(r)
	handleThreadAction(w, r, tm.UnpinThread)
}
