package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	UTORID       = "utorid"
	USER_SESSION = "user-session"
)

var store = sessions.NewCookieStore([]byte("todo-proper-secret")) // TODO: move secret to config

func getSessionUTORid(r *http.Request) (string, bool) {
	session, _ := store.Get(r, USER_SESSION)
	utorid, ok := session.Values[UTORID]
	if !ok {
		return "", ok
	}

	s, ok := utorid.(string)
	return s, ok
}

func isAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := getSessionUTORid(r)

		if ok {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
		}

	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := getSessionUTORid(r); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	utorid := vars["utorid"]

	// TODO: replace this with SAML login for utorid

	session, _ := store.Get(r, USER_SESSION)
	session.Values[UTORID] = utorid
	session.Save(r, w)

	fmt.Fprint(w, "Logged in as: "+utorid)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	utorid, ok := getSessionUTORid(r)
	if !ok {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
	}

	fmt.Fprint(w, "User: "+utorid)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, USER_SESSION)
	delete(session.Values, UTORID)
	session.Options.MaxAge = -1
	session.Save(r, w)
	fmt.Fprint(w, "Logged out")
}
