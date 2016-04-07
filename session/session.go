// Package session provides functionality for user sessions and cookie management.
package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	userSessionName = "user-session"
	userIDKey       = "user-id"
)

// UserSession handles getting and storing user specific properties in a session.
type UserSession struct {
	store sessions.Store
}

// NewUserSession returns a user session.
func NewUserSession(store sessions.Store) *UserSession {
	return &UserSession{store}
}

// getUserSession gets the session containing the user.
func (us *UserSession) get(r *http.Request) (*sessions.Session, error) {
	return us.store.Get(r, userSessionName)
}

// SaveSessionUserID saves the user id in the session. Creates a new session if there was non existing, or overwrites
// if it already exists.
func (us *UserSession) SaveSessionUserID(w http.ResponseWriter, r *http.Request, id int64) error {
	session, err := us.get(r)
	if err != nil {
		return err
	}

	session.Values[userIDKey] = id
	session.Options.HttpOnly = true
	return session.Save(r, w)
}

// SessionUserID gets the id of the user stored in the user session.
// The second return value is a boolean which is true if there is a user id in the session, else false.
func (us *UserSession) SessionUserID(r *http.Request) (int64, bool) {
	session, err := us.get(r)
	if err != nil {
		return -1, false
	}

	i, ok := session.Values[userIDKey]
	if !ok {
		return -1, false
	}

	id, ok := i.(int64)
	return id, ok
}

// Delete deletes the user session.
func (us *UserSession) Delete(w http.ResponseWriter, r *http.Request) error {
	session, err := us.get(r)
	if err != nil {
		return err
	}
	delete(session.Values, userIDKey)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
