package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/httperror"
	"github.com/umairidris/uTeach/models"
)

func addUserUpvotedThreadIDsToData(r *http.Request, threadModel *models.ThreadModel, data map[string]interface{}) error {
	if user, ok := context.SessionUser(r); ok {
		userUpvotedThreadIDs, err := threadModel.GetThreadIdsUpvotedByUser(nil, user)
		if err != nil {
			return err
		}
		data["UserUpvotedThreadIDs"] = userUpvotedThreadIDs
	}
	return nil
}

func getThreads(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)
	data := map[string]interface{}{}
	data["Topic"] = topic

	tm := models.NewThreadModel(a.DB)
	pinnedThreads, err := tm.GetThreadsByTopicAndIsPinned(nil, topic, true)
	if err != nil {
		return err
	}

	unpinnedThreads, err := tm.GetThreadsByTopicAndIsPinned(nil, topic, false)
	if err != nil {
		return err
	}

	data["PinnedThreads"] = pinnedThreads
	data["UnpinnedThreads"] = unpinnedThreads

	if err = addUserUpvotedThreadIDsToData(r, tm, data); err != nil {
		return err
	}

	tagModel := models.NewTagModel(a.DB)
	tags, err := tagModel.GetTagsByTopic(nil, topic)
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
	topicName := strings.ToLower(vars["topic"])
	sm := models.NewTopicModel(a.DB)
	topic, err := sm.GetTopicByName(nil, topicName)
	if err != nil {
		return err
	}

	tm := models.NewTagModel(a.DB)
	tags, err := tm.GetTagsByTopic(nil, topic)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"Tags": tags}
	return renderTemplate(a, w, r, "new_thread.html", data)
}

func postNewThread(a *application.App, w http.ResponseWriter, r *http.Request) error {
	// we want the thread and tags to be created together so use one tx. If one part fails the rest won't be committed.
	tx, err := a.DB.Beginx()
	if err != nil {
		return err
	}

	title := r.FormValue("title")
	text := r.FormValue("text")
	topic := context.Topic(r)
	user, _ := context.SessionUser(r)

	threadModel := models.NewThreadModel(a.DB)
	thread, err := threadModel.AddThread(tx, title, text, topic, user)
	if err != nil {
		return err
	}

	tagIDStr := r.FormValue("tag")
	if tagIDStr != "" {
		tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
		if err != nil {
			return httperror.StatusError{http.StatusBadRequest, err}
		}

		tagModel := models.NewTagModel(a.DB)
		tag, err := tagModel.GetTagByID(nil, tagID)
		if err != nil {
			return err
		}

		if err = tagModel.AddThreadTag(tx, thread, tag); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	http.Redirect(w, r, thread.URL(), http.StatusFound)
	return nil
}

func handleThreadAction(w http.ResponseWriter, r *http.Request, f func(*sqlx.Tx, *models.Thread) error) error {
	thread := context.Thread(r)

	if err := f(nil, thread); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func postThreadVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	tm := models.NewThreadModel(a.DB)

	f := func(tx *sqlx.Tx, thread *models.Thread) error {
		return tm.AddThreadVoteForUser(tx, thread, user)
	}

	return handleThreadAction(w, r, f)
}

func deleteThreadVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	tm := models.NewThreadModel(a.DB)

	f := func(tx *sqlx.Tx, thread *models.Thread) error {
		return tm.RemoveTheadVoteForUser(tx, thread, user)
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

func getThreadsByTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tag := context.Tag(r)

	tm := models.NewThreadModel(a.DB)
	threads, err := tm.GetThreadsByTag(nil, tag)
	if err != nil {
		return err
	}
	data := map[string]interface{}{"Threads": threads}
	if err = addUserUpvotedThreadIDsToData(r, tm, data); err != nil {
		return err
	}
	return renderTemplate(a, w, r, "threads_by_tag.html", data)
}
