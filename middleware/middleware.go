// Package middleware provides app specific middleware handlers.
package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/BrianHarringtonUTSC/uTeach/application"
	"github.com/BrianHarringtonUTSC/uTeach/context"
	"github.com/BrianHarringtonUTSC/uTeach/httperror"
	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/BrianHarringtonUTSC/uTeach/session"
	"github.com/Masterminds/squirrel"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Middleware has app specific middleware.
type Middleware struct {
	App *application.App
}

// SetTemplateData sets the map that contain's template data in the context.
func (m *Middleware) SetTemplateData(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		context.SetTemplateData(r, make(map[string]interface{}))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetSessionUser sets the session user in the context and template data.
func (m *Middleware) SetSessionUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		templateData := context.TemplateData(r)

		us := session.NewUserSession(m.App.Store)
		userID, ok := us.SessionUserID(r)
		if !ok {
			templateData["SessionUser"] = &models.User{}
			next.ServeHTTP(w, r)
			return
		}

		um := models.NewUserModel(m.App.DB)
		user, err := um.FindOne(nil, squirrel.Eq{"users.id": userID})
		if err != nil {
			us.Delete(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		context.SetSessionUser(r, user)
		templateData["SessionUser"] = user
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetTopic sets the topic with the name in the url in the context and template data.
func (m *Middleware) SetTopic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		topicName := strings.ToLower(vars["topicName"])
		tm := models.NewTopicModel(m.App.DB)
		topic, err := tm.FindOne(nil, squirrel.Eq{"topics.name": topicName})
		if err != nil {
			httperror.HandleError(w, errors.Wrap(err, "find one error"))
			return
		}

		context.SetTopic(r, topic)

		templateData := context.TemplateData(r)
		templateData["Topic"] = topic
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetPost sets the post with the id in the url in the context and template data.
func (m *Middleware) SetPost(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID, err := strconv.ParseInt(vars["postID"], 10, 64)
		if err != nil {
			httperror.HandleError(w, httperror.StatusError{http.StatusBadRequest, err})
			return
		}

		pm := models.NewPostModel(m.App.DB)
		topic := context.Topic(r)
		post, err := pm.FindOne(nil, squirrel.Eq{"posts.id": postID, "posts.topic_id": topic.ID})
		if err != nil {
			httperror.HandleError(w, errors.Wrap(err, "find one error"))
			return
		}
		context.SetPost(r, post)

		templateData := context.TemplateData(r)
		templateData["Post"] = post

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetTag sets the tag with name in the url in the context and template data.
func (m *Middleware) SetTag(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tagName := strings.ToLower(vars["tagName"])
		topic := context.Topic(r)

		tm := models.NewTagModel(m.App.DB)
		tag, err := tm.FindOne(nil, squirrel.Eq{"tags.name": tagName, "tags.topic_id": topic.ID})
		if err != nil {
			httperror.HandleError(w, errors.Wrap(err, "find one error"))
			return
		}

		context.SetTag(r, tag)

		templateData := context.TemplateData(r)
		templateData["Tag"] = tag
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// MustLogin ensures the next handler is only accessible by users that are logged in.
func (m *Middleware) MustLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if _, ok := context.SessionUser(r); !ok {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (m *Middleware) isAdmin(r *http.Request) bool {
	user, ok := context.SessionUser(r)
	return ok && user.IsAdmin
}

func (m *Middleware) isPostCreator(r *http.Request) bool {
	post := context.Post(r)
	user, ok := context.SessionUser(r)
	return ok && *post.Creator == *user
}

// MustBeAdmin ensures the next handler is only accessible by an admin.
func (m *Middleware) MustBeAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !m.isAdmin(r) {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// MustBeAdminOrPostCreator ensures the next handler is only accessible by an admin or the post creator.
func (m *Middleware) MustBeAdminOrPostCreator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !m.isPostCreator(r) && !m.isAdmin(r) {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
