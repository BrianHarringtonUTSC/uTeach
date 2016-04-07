// Package context provides convenience setters and getters for request-scoped context values.
package context

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/umairidris/uTeach/models"
)

// TODO: Use context in standard lib which will be added in go 1.7

const (
	subjectKey     = "subject"
	threadKey      = "thread"
	tagKey         = "tag"
	sessionUserKey = "session-user"
)

// SetSubject sets the subject in the context.
func SetSubject(r *http.Request, subject *models.Subject) {
	context.Set(r, subjectKey, subject)
}

// Subject gets the subject from the context.
func Subject(r *http.Request) *models.Subject {
	return context.Get(r, subjectKey).(*models.Subject)
}

// SetThread sets the thread in the context.
func SetThread(r *http.Request, thread *models.Thread) {
	context.Set(r, threadKey, thread)
}

// Thread gets the thread from the context.
func Thread(r *http.Request) *models.Thread {
	return context.Get(r, threadKey).(*models.Thread)
}

// SetTag sets the tag in the context.
func SetTag(r *http.Request, thread *models.Tag) {
	context.Set(r, tagKey, thread)
}

// Tag sets the tag in the context.
func Tag(r *http.Request) *models.Tag {
	return context.Get(r, tagKey).(*models.Tag)
}

// SetSessionUser sets the session user in the context.
func SetSessionUser(r *http.Request, thread *models.User) {
	context.Set(r, sessionUserKey, thread)
}

// SessionUser gets the session user from the context.
func SessionUser(r *http.Request) (*models.User, bool) {
	user, ok := context.Get(r, sessionUserKey).(*models.User)
	return user, ok
}
