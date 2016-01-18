package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	USER_KEY     = "user"
	USER_SESSION = "user-session"
)

var store = sessions.NewCookieStore([]byte("todo-proper-secret")) // TODO: move secret to config

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, USER_SESSION)
}

func getSessionUser(r *http.Request) (*User, bool) {
	session, err := GetSession(r)
	if err != nil {
		return nil, false
	}

	user, ok := session.Values[USER_KEY]
	if !ok {
		return nil, ok
	}

	u, ok := user.(*User)
	return u, ok
}

func isAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := getSessionUser(r)

		if ok {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
		}
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := getSessionUser(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	utorid := vars["utorid"]

	// TODO: replace this with SAML login for utorid

	session, err := GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, ok := GetUser(utorid)
	if !ok {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	session.Values[USER_KEY] = user
	session.Save(r, w)

	fmt.Fprint(w, "Logged in as: "+utorid)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	user, ok := getSessionUser(r)
	if !ok {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User: "+user.UTORid)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(session.Values, USER_KEY)
	session.Options.MaxAge = -1
	session.Save(r, w)
	fmt.Fprint(w, "Logged out")
}
