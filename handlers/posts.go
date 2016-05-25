package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/BrianHarringtonUTSC/uTeach/application"
	"github.com/BrianHarringtonUTSC/uTeach/context"
	"github.com/BrianHarringtonUTSC/uTeach/httperror"
	"github.com/BrianHarringtonUTSC/uTeach/libtemplate"
	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func addUserUpvotedPostIDsToData(r *http.Request, postModel *models.PostModel, data map[string]interface{}) error {
	if user, ok := context.SessionUser(r); ok {
		userUpvotedPostIDs, err := postModel.GetVotedPostIds(nil, squirrel.Eq{"post_votes.user_id": user.ID})
		if err != nil {
			return errors.Wrap(err, "get upvoted post ids error")
		}
		data["UserUpvotedPostIDs"] = userUpvotedPostIDs
	}
	return nil
}

func getPosts(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)

	whereEq := squirrel.Eq{"posts.topic_id": topic.ID, "posts.is_pinned": true}

	pm := models.NewPostModel(a.DB)
	pinnedPosts, err := pm.Find(nil, whereEq)
	switch {
	case err == sql.ErrNoRows:
		pinnedPosts = make([]*models.Post, 0)
	case err != nil:
		return errors.Wrap(err, "find error")
	}

	whereEq["posts.is_pinned"] = false
	unpinnedPosts, err := pm.Find(nil, whereEq)
	switch {
	case err == sql.ErrNoRows:
		unpinnedPosts = make([]*models.Post, 0)
	case err != nil:
		return errors.Wrap(err, "find error")
	}

	tagModel := models.NewTagModel(a.DB)
	tags, err := tagModel.Find(nil, squirrel.Eq{"tags.topic_id": topic.ID})
	if err != nil {
		return errors.Wrap(err, "find error")
	}

	data := context.TemplateData(r)
	data["PinnedPosts"] = pinnedPosts
	data["UnpinnedPosts"] = unpinnedPosts
	data["Tags"] = tags

	if err = addUserUpvotedPostIDsToData(r, pm, data); err != nil {
		return errors.Wrap(err, "add upvoted post ids to data error")
	}

	return libtemplate.Render(w, a.Templates, "posts.html", data)
}

func getPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return libtemplate.Render(w, a.Templates, "post.html", context.TemplateData(r))
}

func getNewPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)

	tm := models.NewTagModel(a.DB)
	tags, err := tm.Find(nil, squirrel.Eq{"tags.topic_id": topic.ID})
	if err != nil {
		return errors.Wrap(err, "find error")
	}

	data := context.TemplateData(r)
	data["Tags"] = tags
	return libtemplate.Render(w, a.Templates, "new_post.html", data)
}

func postNewPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	title := r.FormValue("title")
	text := r.FormValue("text")
	topic := context.Topic(r)
	user, _ := context.SessionUser(r)

	// we want the post and tags to be created together so use one tx. If one part fails the rest won't be committed.
	tx, err := a.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "begin transacion error")
	}

	postModel := models.NewPostModel(a.DB)
	post, err := postModel.AddPost(tx, title, text, topic, user)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "add post error")
	}

	tagIDStr := r.FormValue("tag")
	if tagIDStr != "" {
		tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
		if err != nil {
			tx.Rollback()
			return httperror.StatusError{http.StatusBadRequest, err}
		}

		tagModel := models.NewTagModel(a.DB)
		tag, err := tagModel.FindOne(nil, squirrel.Eq{"tags.id": tagID})
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "find one error")
		}

		if err = tagModel.AddPostTag(tx, post, tag); err != nil {
			tx.Rollback()
			return errors.Wrap(err, "add post tag error")
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "commit error")
	}

	http.Redirect(w, r, post.URL(), http.StatusFound)
	return nil
}

func handlePostAction(w http.ResponseWriter, r *http.Request, f func(*sqlx.Tx, *models.Post) error) error {
	post := context.Post(r)

	if err := f(nil, post); err != nil {
		return errors.Wrap(err, "post error")
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func postPostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	pm := models.NewPostModel(a.DB)

	f := func(tx *sqlx.Tx, post *models.Post) error {
		err := pm.AddPostVoteForUser(tx, post, user)
		return errors.Wrap(err, "add post vote error")
	}

	err := handlePostAction(w, r, f)
	return errors.Wrap(err, "handle post action error")
}

func deletePostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	pm := models.NewPostModel(a.DB)

	f := func(tx *sqlx.Tx, post *models.Post) error {
		err := pm.RemovePostVoteForUser(tx, post, user)
		return errors.Wrap(err, "remove post vote error")
	}

	err := handlePostAction(w, r, f)
	return errors.Wrap(err, "handle post action error")
}

func postHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	err := handlePostAction(w, r, pm.HidePost)
	return errors.Wrap(err, "hide post error")
}

func deleteHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	err := handlePostAction(w, r, pm.UnhidePost)
	return errors.Wrap(err, "unhide post error")
}

func postPinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	err := handlePostAction(w, r, pm.PinPost)
	return errors.Wrap(err, "pin post error")
}

func deletePinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	err := handlePostAction(w, r, pm.UnpinPost)
	return errors.Wrap(err, "unpin post error")
}

func getPostsByTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tag := context.Tag(r)

	pm := models.NewPostModel(a.DB)
	posts, err := pm.Find(nil, squirrel.Eq{"post_tags.tag_id": tag.ID})
	if err != nil {
		return errors.Wrap(err, "find error")
	}

	data := context.TemplateData(r)
	data["Posts"] = posts
	if err = addUserUpvotedPostIDsToData(r, pm, data); err != nil {
		return errors.Wrap(err, "add upvoted post ids to data error")
	}

	err = libtemplate.Render(w, a.Templates, "posts_by_tag.html", data)
	return errors.Wrap(err, "render template error")
}
