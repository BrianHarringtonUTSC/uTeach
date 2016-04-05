package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
)

func getThreads(a *application.App, w http.ResponseWriter, r *http.Request) error {
	subject := context.Subject(r)

	tm := models.NewThreadModel(a.DB)
	pinnedThreads, err := tm.GetThreadsBySubjectAndIsPinned(nil, subject, true)
	if err != nil {
		return err
	}

	unpinnedThreads, err := tm.GetThreadsBySubjectAndIsPinned(nil, subject, false)
	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	data["PinnedThreads"] = pinnedThreads
	data["UnpinnedThreads"] = unpinnedThreads

	//  if there is a user, get the user's upvoted threads
	if user, ok := context.SessionUser(r); ok {
		userUpvotedThreadIDs, err := tm.GetThreadIdsUpvotedByUser(nil, user)
		if err != nil {
			return err
		}
		data["UserUpvotedThreadIDs"] = userUpvotedThreadIDs
	}

	tagModel := models.NewTagModel(a.DB)
	tags, err := tagModel.GetTagsBySubject(nil, subject)
	if err != nil {
		return err
	}

	data["Tags"] = tags

	return renderTemplate(a, w, r, "threads.html", data)
}

func getThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	thread := context.Thread(r)
	data := map[string]interface{}{"Thread": thread}
	return renderTemplate(a, w, r, "thread.html", data)
}

func getNewThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subjectName := strings.ToLower(vars["subject"])
	sm := models.NewSubjectModel(a.DB)
	subject, err := sm.GetSubjectByName(nil, subjectName)
	if err != nil {
		return err
	}

	tm := models.NewTagModel(a.DB)
	tags, err := tm.GetTagsBySubject(nil, subject)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Tags": tags}
	return renderTemplate(a, w, r, "new_thread.html", data)
}

func postNewThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	subject := context.Subject(r)

	user, _ := context.SessionUser(r)

	title := r.FormValue("title")
	text := r.FormValue("text")

	tx, err := a.DB.Beginx()
	if err != nil {
		return err
	}
	tm := models.NewThreadModel(a.DB)
	thread, err := tm.AddThread(tx, title, text, subject, user)
	if err != nil {
		return err
	}

	tagIDStr := r.FormValue("tag")
	if tagIDStr != "" {
		tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
		if err != nil {
			return err
		}
		tagModel := models.NewTagModel(a.DB)
		tag, err := tagModel.GetTagByID(nil, tagID)
		if err != nil {
			return err
		}
		err = tagModel.AddThreadTag(tx, thread, tag)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	http.Redirect(w, r, thread.URL(), 301)
	return nil
}

func handleThreadAction(w http.ResponseWriter, r *http.Request, f func(*sqlx.Tx, int64) error) error {
	thread := context.Thread(r)

	if err := f(nil, thread.ID); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func postThreadVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	tm := models.NewThreadModel(a.DB)

	f := func(tx *sqlx.Tx, id int64) error {
		return tm.AddThreadVoteForUser(tx, id, user)
	}

	return handleThreadAction(w, r, f)
}

func deleteThreadVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	tm := models.NewThreadModel(a.DB)

	f := func(tx *sqlx.Tx, id int64) error {
		return tm.RemoveTheadVoteForUser(tx, id, user)
	}

	return handleThreadAction(w, r, f)
}

func postHideThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.HideThread)
}

func deleteHideThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.UnhideThread)
}

func postPinThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.PinThread)
}

func deletePinThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewThreadModel(a.DB)
	return handleThreadAction(w, r, tm.UnpinThread)
}
