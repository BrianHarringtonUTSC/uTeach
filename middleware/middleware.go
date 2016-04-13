// Package middleware provides app specific middleware handlers.
package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/httperror"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

// Middleware has app specific middleware.
type Middleware struct {
	App *application.App
}

// SetSessionUser sets the session user in the context.
func (m *Middleware) SetSessionUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		us := session.NewUserSession(m.App.Store)
		userID, ok := us.SessionUserID(r)
		if ok {
			um := models.NewUserModel(m.App.DB)
			user, err := um.GetUserByID(nil, userID)
			if err != nil {
				us.Delete(w, r)
				http.Redirect(w, r, "/", http.StatusFound)
			}
			context.SetSessionUser(r, user)
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetTopic sets the topic with the name in the url in the context.
func (m *Middleware) SetTopic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		topicName := strings.ToLower(vars["topic"])
		sm := models.NewTopicModel(m.App.DB)
		topic, err := sm.GetTopicByName(nil, topicName)
		if err != nil {
			httperror.HandleError(w, err)
			return
		}

		context.SetTopic(r, topic)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetThread sets the thread with the id in the url in the context.
func (m *Middleware) SetThread(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
		if err != nil {
			httperror.HandleError(w, httperror.StatusError{http.StatusBadRequest, err})
			return
		}

		tm := models.NewThreadModel(m.App.DB)
		thread, err := tm.GetThreadByID(nil, threadID)
		if err != nil {
			httperror.HandleError(w, err)
			return
		}
		context.SetThread(r, thread)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SetTag sets the tag with name in the url in the context.
func (m *Middleware) SetTag(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tagName := strings.ToLower(vars["tag"])
		tm := models.NewTagModel(m.App.DB)
		tag, err := tm.GetTagByNameAndTopic(nil, tagName, context.Topic(r))
		if err != nil {
			httperror.HandleError(w, err)
			return
		}

		context.SetTag(r, tag)
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

func (m *Middleware) isThreadCreator(r *http.Request) bool {
	thread := context.Thread(r)
	user, ok := context.SessionUser(r)
	return ok && thread.Creator == user
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

// MustBeAdminOrThreadCreator ensures the next handler is only accessible by an admin or the thread creator.
func (m *Middleware) MustBeAdminOrThreadCreator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !m.isThreadCreator(r) && !m.isAdmin(r) {
			httperror.HandleError(w, httperror.StatusError{http.StatusForbidden, nil})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
