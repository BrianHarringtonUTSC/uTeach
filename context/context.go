// Package context provides convenience setters and getters for request-scoped context values.
package context

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/umairidris/uTeach/models"
)

// TODO: replace this with context which will be added in go 1.7
const (
	subjectKey     = "subject"
	threadKey      = "thread"
	tagKey         = "tag"
	sessionUserKey = "session-user"
)

func SetSubject(r *http.Request, subject *models.Subject) {
	context.Set(r, subjectKey, subject)
}

func Subject(r *http.Request) *models.Subject {
	return context.Get(r, subjectKey).(*models.Subject)
}

func SetThread(r *http.Request, thread *models.Thread) {
	context.Set(r, threadKey, thread)
}

func Thread(r *http.Request) *models.Thread {
	return context.Get(r, threadKey).(*models.Thread)
}

func SetTag(r *http.Request, thread *models.Tag) {
	context.Set(r, tagKey, thread)
}

func Tag(r *http.Request) *models.Tag {
	return context.Get(r, tagKey).(*models.Tag)
}

func SetSessionUser(r *http.Request, thread *models.User) {
	context.Set(r, sessionUserKey, thread)
}

func SessionUser(r *http.Request) (*models.User, bool) {
	user, ok := context.Get(r, sessionUserKey).(*models.User)
	return user, ok
}
