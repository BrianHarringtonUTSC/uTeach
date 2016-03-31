package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

// GetThreads renders all threads for the subject.
func GetThreads(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subject_name := strings.ToLower(vars["subject"])

	sm := models.NewSubjectModel(a.DB)
	subject, err := sm.GetSubjectByName(subject_name)
	if err != nil {
		return err
	}

	tm := models.NewThreadModel(a.DB)
	pinnedThreads, err := tm.GetThreadsBySubjectAndIsPinned(subject, true)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	unpinnedThreads, err := tm.GetThreadsBySubjectAndIsPinned(subject, false)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	data := map[string]interface{}{}
	data["PinnedThreads"] = pinnedThreads
	data["UnpinnedThreads"] = unpinnedThreads

	//  if there is a user, get the user's upvoted threads
	usm := session.NewUserSessionManager(a.CookieStore)
	if user, ok := usm.SessionUser(r); ok {
		userUpvotedThreadIDs, err := tm.GetThreadIdsUpvotedByUser(user)
		if err != nil {
			return err
		}
		data["UserUpvotedThreadIDs"] = userUpvotedThreadIDs
	}

	return renderTemplate(a, w, r, "threads.html", data)
}

// GetThread renders a thread.
func GetThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)

	threadID := context.ThreadID(r)
	thread, err := tm.GetThreadByID(threadID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Thread": thread}
	return renderTemplate(a, w, r, "thread.html", data)
}

// GetNewThread renders the new thread page.
func GetNewThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(a, w, r, "new_thread.html", nil)
}

// PostNewThread adds a new thread in the db and redirects to it, if successful.
func PostNewThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subject_name := strings.ToLower(vars["subject"])

	sm := models.NewSubjectModel(a.DB)
	subject, err := sm.GetSubjectByName(subject_name)
	if err != nil {
		return err
	}

	usm := session.NewUserSessionManager(a.CookieStore)
	user, _ := usm.SessionUser(r)

	title := r.FormValue("title")
	text := r.FormValue("text")

	tm := models.NewThreadModel(a.DB)
	thread, err := tm.AddThread(title, text, subject, user)
	if err != nil {
		return err
	}
	http.Redirect(w, r, thread.URL(), 301)
	return nil
}

func handleThreadAction(w http.ResponseWriter, r *http.Request, f func(int64) error) error {
	threadID := context.ThreadID(r)

	if err := f(threadID); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// DeleteThreadVote removes a vote for the user on a thread.
func PostThreadVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	usm := session.NewUserSessionManager(a.CookieStore)
	user, _ := usm.SessionUser(r)

	tm := models.NewThreadModel(a.DB)

	f := func(id int64) error {
		return tm.AddThreadVoteForUser(id, user)
	}

	return handleThreadAction(w, r, f)
}

// DeleteThreadVote removes a vote for the user on a thread.
func DeleteThreadVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	usm := session.NewUserSessionManager(a.CookieStore)
	user, _ := usm.SessionUser(r)

	tm := models.NewThreadModel(a.DB)

	f := func(id int64) error {
		return tm.RemoveTheadVoteForUser(id, user)
	}

	return handleThreadAction(w, r, f)
}

func PostHideThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.HideThread)
}

func DeleteHideThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.UnhideThread)
}

func PostPinThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.PinThread)
}

func DeletePinThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.UnpinThread)
}
