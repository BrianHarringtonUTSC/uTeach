package handlers

import (
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/httperror"
	"github.com/umairidris/uTeach/libtemplate"
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
	pm := models.NewPostModel(a.DB)
	pinnedPosts, err := pm.GetPostsByTopicAndIsPinned(nil, topic, true)
	if err != nil {
		return err
	}

	unpinnedPosts, err := pm.GetPostsByTopicAndIsPinned(nil, topic, false)
	if err != nil {
		return err
	}
	tagModel := models.NewTagModel(a.DB)
	tags, err := tagModel.GetTagsByTopic(nil, topic)
	if err != nil {
		return err
	}

	data := context.TemplateData(r)
	data["PinnedPosts"] = pinnedPosts
	data["UnpinnedPosts"] = unpinnedPosts
	data["Tags"] = tags

	if err = addUserUpvotedPostIDsToData(r, pm, data); err != nil {
		return err
	}

	return libtemplate.Render(w, a.Templates, "posts.html", data)
}

func getPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	return libtemplate.Render(w, a.Templates, "post.html", context.TemplateData(r))
}

func getNewPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	topic := context.Topic(r)

	tm := models.NewTagModel(a.DB)
	tags, err := tm.GetTagsByTopic(nil, topic)
	if err != nil {
		return err
	}

	data := context.TemplateData(r)
	data["Tags"] = tags
	return libtemplate.Render(w, a.Templates, "new_post.html", data)
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

	pm := models.NewPostModel(a.DB)

	f := func(tx *sqlx.Tx, post *models.Post) error {
		return pm.AddPostVoteForUser(tx, post, user)
	}

	return handlePostAction(w, r, f)
}

func deletePostVote(a *application.App, w http.ResponseWriter, r *http.Request) error {
	user, _ := context.SessionUser(r)

	pm := models.NewPostModel(a.DB)

	f := func(tx *sqlx.Tx, post *models.Post) error {
		return pm.RemoveTheadVoteForUser(tx, post, user)
	}

	return handlePostAction(w, r, f)
}

func postHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, pm.HidePost)
}

func deleteHidePost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, pm.UnhidePost)
}

func postPinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, pm.PinPost)
}

func deletePinPost(a *application.App, w http.ResponseWriter, r *http.Request) error {
	pm := models.NewPostModel(a.DB)
	return handlePostAction(w, r, pm.UnpinPost)
}

func getPostsByTag(a *application.App, w http.ResponseWriter, r *http.Request) error {
	tag := context.Tag(r)

	pm := models.NewPostModel(a.DB)
	posts, err := pm.GetPostsByTag(nil, tag)
	if err != nil {
		return err
	}

	data := context.TemplateData(r)
	data["Posts"] = posts
	if err = addUserUpvotedPostIDsToData(r, pm, data); err != nil {
		return err
	}
	return libtemplate.Render(w, a.Templates, "posts_by_tag.html", data)
}
