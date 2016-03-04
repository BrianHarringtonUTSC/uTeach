// Package middleware provides middleware handlers for the uTeach app.
package middleware

import (
	"net/http"

	"github.com/umairidris/uTeach/application"
)

// SetApplication sets the application in the context for other handlers to use.
func SetApplication(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			application.SetInContext(r, app)
			next.ServeHTTP(w, r)
		})
	}
}

// MustLogin ensures that the next handler is only accessible by users that are logged in.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app := application.GetFromContext(r)
		if _, ok := app.Store.SessionUser(r); !ok {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isSessionUserAdmin(r *http.Request) bool {
	app := application.GetFromContext(r)
	user, _ := app.Store.SessionUser(r)
	return user.IsAdmin
}

func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !isSessionUserAdmin(r) {
			http.Error(w, "You must be an admin to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MustBeAdminOrThreadOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO
		next.ServeHTTP(w, r)
	})
}

func SetThreadID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO
		next.ServeHTTP(w, r)
	})
}
