package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	UserKey         = "user"
	UserSessionName = "user-session"
)

type Store struct {
	*sessions.CookieStore
}

func NewStore(authenticationKey string) *Store {
	return &Store{sessions.NewCookieStore([]byte(authenticationKey))}
}

func (s *Store) getUserSession(r *http.Request) (*sessions.Session, error) {
	return s.Get(r, UserSessionName)
}

func (s *Store) SessionUser(r *http.Request) (*User, bool) {
	session, err := s.getUserSession(r)
	if err != nil {
		return nil, false
	}

	user, ok := session.Values[UserKey]
	if !ok {
		return nil, false
	}

	u, ok := user.(*User)
	return u, ok
}

func (s *Store) NewUserSession(w http.ResponseWriter, r *http.Request, username string, db *DB) error {
	session, err := s.getUserSession(r)
	if err != nil {
		return err
	}

	user, err := db.GetUser(username)
	if err != nil {
		return err
	}

	session.Values[UserKey] = user
	return session.Save(r, w)
}

func (s *Store) DeleteUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.getUserSession(r)
	if err != nil {
		return err
	}
	delete(session.Values, UserKey)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
