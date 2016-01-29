package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	USER_KEY     = "user"
	USER_SESSION = "user-session"
)

var store = sessions.NewCookieStore([]byte("todo-proper-secret")) // TODO: move secret to config

func GetUserSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, USER_SESSION)
}

func GetSessionUser(r *http.Request) (*User, bool) {
	session, err := GetUserSession(r)
	if err != nil {
		return nil, false
	}

	user, ok := session.Values[USER_KEY]
	if !ok {
		return nil, false
	}

	u, ok := user.(*User)
	return u, ok
}

func NewUserSession(w http.ResponseWriter, r *http.Request, username string) error {
	session, err := GetUserSession(r)
	if err != nil {
		return err
	}

	user, err := GetUser(username)
	if err != nil {
		return err
	}

	session.Values[USER_KEY] = user
	return session.Save(r, w)
}

func DeleteUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := GetUserSession(r)
	if err != nil {
		return err
	}
	delete(session.Values, USER_KEY)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
