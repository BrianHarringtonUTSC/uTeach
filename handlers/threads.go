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

func addUserUpvotedPostIDsToData(r *http.Request, postModel *models.PostModel, data map[string]interface{}) error {
	if user, ok := context.SessionUser(r); ok {
		userUpvotedPostIDs, err := postModel.GetPostIdsUpvotedByUser(nil, user)
		if err != nil {
			return err
		}
		data["UserUpvotedPostIDs"] = userUpvotedPostIDs
	}
	return nil
}

func getPosts(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)
	data := map[string]interface{}{}
	data["Topic"] = topic

	tm := models.NewPostModel(a.DB)
	pinnedPosts, err := tm.GetPostsByTopicAndIsPinned(nil, topic, true)
	if err != nil {
		return err
	}

	unpinnedPosts, err := tm.GetPostsByTopicAndIsPinned(nil, topic, false)
	if err != nil {
		return err
	}

	data["PinnedPosts"] = pinnedPosts
	data["UnpinnedPosts"] = unpinnedPosts

	if err = addUserUpvotedPostIDsToData(r, tm, data); err != nil {
		return err
	}

	tagModel := models.NewTagModel(a.DB)
	tags, err := tagModel.GetTagsByTopic(nil, topic)
	if err != nil {
		return err
	}

	data["Tags"] = tags

	return renderTemplate(a, w, r, "posts.html", data)
}

func getPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	post := context.Post(r)
	data := map[string]interface{}{"Post": post}
	return renderTemplate(a, w, r, "post.html", data)
}

func getNewPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
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
	return renderTemplate(a, w, r, "new_post.html", data)
}

func postNewPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	// we want the post and tags to be created together so use one tx. If one part fails the rest won't be committed.
	tx, err := a.DB.Beginx()
	if err != nil {
		return err
	}

	title := r.FormValue("title")
	text := r.FormValue("text")
	topic := context.Topic(r)
	user, _ := context.SessionUser(r)

	postModel := models.NewPostModel(a.DB)
	post, err := postModel.AddPost(tx, title, text, topic, user)
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

		if err = tagModel.AddPostTag(tx, post, tag); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	http.Redirect(w, r, post.URL(), http.StatusFound)
	return nil
}

func handlePostAction(w http.ResponseWriter, r *http.Request, f func(*sqlx.Tx, *models.Post) error) error {
	post := context.Post(r)

	if err := f(nil, post); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func postPostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	tm := models.NewPostModel(a.DB)

	f := func(tx *sqlx.Tx, post *models.Post) error {
		return tm.AddPostVoteForUser(tx, post, user)
	}

	return handlePostAction(w, r, f)
}

func deletePostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	tm := models.NewPostModel(a.DB)

	f := func(tx *sqlx.Tx, post *models.Post) error {
		return tm.RemoveTheadVoteForUser(tx, post, user)
	}

	return handlePostAction(w, r, f)
}

func postHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, tm.HidePost)
}

func deleteHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, tm.UnhidePost)
}

func postPinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, tm.PinPost)
}

func deletePinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, tm.UnpinPost)
}

func getPostsByTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tag := context.Tag(r)

	tm := models.NewPostModel(a.DB)
	posts, err := tm.GetPostsByTag(nil, tag)
	if err != nil {
		return err
	}
	data := map[string]interface{}{"Posts": posts}
	if err = addUserUpvotedPostIDsToData(r, tm, data); err != nil {
		return err
	}
	return renderTemplate(a, w, r, "posts_by_tag.html", data)
}
