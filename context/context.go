// Package context provides convenience setters and getters for request-scoped context values.
package context

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/umairidris/uTeach/models"
)

// TODO: Use context in standard lib which will be added in go 1.7

const (
	topicKey       = "topic"
	postKey        = "post"
	tagKey         = "tag"
	sessionUserKey = "session-user"
)

// SetTopic sets the topic in the context.
func SetTopic(r *http.Request, topic *models.Topic) {
	context.Set(r, topicKey, topic)
}

// Topic gets the topic from the context.
func Topic(r *http.Request) *models.Topic {
	return context.Get(r, topicKey).(*models.Topic)
}

// SetPost sets the post in the context.
func SetPost(r *http.Request, post *models.Post) {
	context.Set(r, postKey, post)
}

// Post gets the post from the context.
func Post(r *http.Request) *models.Post {
	return context.Get(r, postKey).(*models.Post)
}

// SetTag sets the tag in the context.
func SetTag(r *http.Request, post *models.Tag) {
	context.Set(r, tagKey, post)
}

// Tag sets the tag in the context.
func Tag(r *http.Request) *models.Tag {
	return context.Get(r, tagKey).(*models.Tag)
}

// SetSessionUser sets the session user in the context.
func SetSessionUser(r *http.Request, post *models.User) {
	context.Set(r, sessionUserKey, post)
}

// SessionUser gets the session user from the context.
func SessionUser(r *http.Request) (*models.User, bool) {
	user, ok := context.Get(r, sessionUserKey).(*models.User)
	return user, ok
}
