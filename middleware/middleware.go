// Package middleware provides middleware handlers for the uTeach app.
package middleware

import (
	"net/http"

	"github.com/umairidris/uTeach/app"
)

type Middleware struct {
	App *app.App
}

// AuthorizedHandler ensures that the next handler is only accessible by users that are logged in.
func (m *Middleware) AuthorizedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if _, ok := m.App.Store.SessionUser(r); !ok {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
