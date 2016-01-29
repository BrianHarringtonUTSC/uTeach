package main

import (
	"net/http"
)

func AuthorizedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetSessionUser(r)

		if ok {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
		}
	})
}
