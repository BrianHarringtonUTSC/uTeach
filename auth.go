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

// session store
var store = sessions.NewCookieStore([]byte("todo-proper-secret")) // TODO: move secret to config

func getSessionUTORid(session *sessions.Session) (string, bool) {
	utorid, ok := session.Values[UTORID]
	if !ok {
		return "", ok
	}

	s, ok := utorid.(string)
	return s, ok
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, USER_SESSION)
	if _, ok := getSessionUTORid(session); ok {
		fmt.Fprint(w, "Already logged in")
		return
	}

	vars := mux.Vars(r)
	utorid := vars["utorid"]

	// TODO: replace this with SAML login for utorid

	session.Values[UTORID] = utorid
	session.Save(r, w)

	fmt.Fprint(w, "Logged in as: "+utorid)
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, USER_SESSION)
	utorid, ok := getSessionUTORid(session)
	msg := "Not logged in"
	if ok {
		msg = "logged in as: " + utorid
	}

	fmt.Fprint(w, msg)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, USER_SESSION)
	session.Values[UTORID] = nil
	session.Save(r, w)
	fmt.Fprint(w, "Logged out")
}
