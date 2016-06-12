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

func postNewPost(a *application.App, w http.ResponseWriter, r *http.Request) (err error) {
	title := r.FormValue("title")
	text := r.FormValue("text")
	topic := context.Topic(r)
	user, _ := context.SessionUser(r)

	// we want the post and tags to be created together so use one tx. If one part fails the rest won't be committed.
	tx, err := a.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "begin transacion error")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		err = errors.Wrap(err, "commit error")
	}()

	postModel := models.NewPostModel(a.DB)
	post := &models.Post{Title: title, Content: text, Topic: topic, Creator: user}
	if err = postModel.Add(tx, post); err != nil {
		return errors.Wrap(err, "add post error")
	}

	tagIDStr := r.FormValue("tag")
	if tagIDStr != "" {
		tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
		if err != nil {
			return httperror.StatusError{http.StatusBadRequest, err}
		}

		tagModel := models.NewTagModel(a.DB)
		tag, err := tagModel.FindOne(nil, squirrel.Eq{"tags.id": tagID})
		if err != nil {
			return errors.Wrap(err, "find one error")
		}

		if err = tagModel.AddPostTag(tx, post, tag); err != nil {
			return errors.Wrap(err, "add post tag error")
		}
	}

	http.Redirect(w, r, post.URL(), http.StatusFound)
	return nil
}

func postHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	post := context.Post(r)
	post.IsVisible = false
	err := pm.Update(nil, post)
	return errors.Wrap(err, "update error")
}

func deleteHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	post := context.Post(r)
	post.IsVisible = true
	err := pm.Update(nil, post)
	return errors.Wrap(err, "update error")
}

func postPinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	post := context.Post(r)
	post.IsPinned = true
	err := pm.Update(nil, post)
	return errors.Wrap(err, "update error")
}

func deletePinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	post := context.Post(r)
	post.IsPinned = false
	err := pm.Update(nil, post)
	return errors.Wrap(err, "update error")
}

func updatePostVote(a *application.App, w http.ResponseWriter, r *http.Request, voted bool) error {
	post := context.Post(r)
	user, _ := context.SessionUser(r)
	pm := models.NewPostModel(a.DB)

	if err := pm.UpdatePostVoteForUser(nil, post, user, voted); err != nil {
		return errors.Wrap(err, "update post vote error")
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func postPostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return updatePostVote(a, w, r, true)
}

func deletePostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return updatePostVote(a, w, r, false)
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
