package main

import (
	"net/http"
)

func (a *App) AuthorizedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, ok := a.store.SessionUser(r); !ok {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
