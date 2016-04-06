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

type UserSession struct {
	store sessions.Store
}

func NewUserSession(store sessions.Store) *UserSession {
	return &UserSession{store}
}

// getUserSession gets the session containing the user.
func (us *UserSession) get(r *http.Request) (*sessions.Session, error) {
	return us.store.Get(r, userSessionName)
}

// New creates a new session and stores the User containing the
func (us *UserSession) SaveSessionUserID(w http.ResponseWriter, r *http.Request, id int64) error {
	session, err := us.get(r)
	if err != nil {
		return err
	}

	session.Values[userIDKey] = id
	session.Options.HttpOnly = true
	return session.Save(r, w)
}

// User gets the user stored in the user session.
// If there is a User stored in the session and can be retrieved it returns the user and true, else the boolean will be
// false.
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
