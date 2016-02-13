package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	userSessionName = "user-session"
	userKey         = "user"
)

// Store wraps around CookieStore to provide uTeach specific cookie and session functionality.
type Store struct {
	*sessions.CookieStore
}

// NewStore creates a new store.
func NewStore(authenticationKey string) *Store {
	return &Store{sessions.NewCookieStore([]byte(authenticationKey))}
}

// getUserSession gets the session containing the user.
func (s *Store) getUserSession(r *http.Request) (*sessions.Session, error) {
	return s.Get(r, userSessionName)
}

// NewUserSession creates a new session and stores the User containing the
func (s *Store) NewUserSession(w http.ResponseWriter, r *http.Request, username string, db *DB) error {
	session, err := s.getUserSession(r)
	if err != nil {
		return err
	}

	user, err := db.User(username)
	if err != nil {
		return err
	}

	session.Values[userKey] = user
	return session.Save(r, w)
}

// SessionUser gets the user stored in the user session.
// If there is a User stored in the session and can be retrieved it returns the user and true, else the boolean will be
// false.
func (s *Store) SessionUser(r *http.Request) (user *User, ok bool) {
	session, err := s.getUserSession(r)
	if err != nil {
		return nil, false
	}

	u, ok := session.Values[userKey]
	if !ok {
		return nil, false
	}

	user, ok = u.(*User)
	return
}

// DeleteUserSession deletes the user session.
func (s *Store) DeleteUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.getUserSession(r)
	if err != nil {
		return err
	}
	delete(session.Values, userKey)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
