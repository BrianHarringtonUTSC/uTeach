// Package context provides convenience setters and getters for request-scoped context values.
package context

import (
	"net/http"

	"github.com/BrianHarringtonUTSC/uTeach/models"
	"github.com/gorilla/context"
)

// TODO: Use context in standard lib which will be added in go 1.7

const (
	templateDataKey = "template-data"
	topicKey        = "topic"
	postKey         = "post"
	tagKey          = "tag"
	sessionUserKey  = "session-user"
)

// SetTemplateData sets the template data map in the context.
func SetTemplateData(r *http.Request, data map[string]interface{}) {
	context.Set(r, templateDataKey, data)
}

// TemplateData gets the template data map from the context.
func TemplateData(r *http.Request) map[string]interface{} {
	return context.Get(r, templateDataKey).(map[string]interface{})
}

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
func SetTag(r *http.Request, tag *models.Tag) {
	context.Set(r, tagKey, tag)
}

// Tag sets the tag in the context.
func Tag(r *http.Request) *models.Tag {
	return context.Get(r, tagKey).(*models.Tag)
}

// SetSessionUser sets the session user in the context.
func SetSessionUser(r *http.Request, user *models.User) {
	context.Set(r, sessionUserKey, user)
}

// SessionUser gets the session user from the context.
func SessionUser(r *http.Request) (*models.User, bool) {
	user, ok := context.Get(r, sessionUserKey).(*models.User)
	return user, ok
}
