// Package middleware provides middleware handlers for the uTeach app.
package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/context"
	"github.com/umairidris/uTeach/models"
	"github.com/umairidris/uTeach/session"
)

// SetApplication sets the application in the context for other handlers to use.
func SetApplication(a *application.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.SetApp(r, a)
			next.ServeHTTP(w, r)
		})
	}
}

func SetThreadIDVar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		threadID, err := strconv.ParseInt(vars["threadID"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		context.SetThreadID(r, threadID)
		next.ServeHTTP(w, r)
	})

}

// MustLogin ensures that the next handler is only accessible by users that are logged in.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		usm := session.NewUserSessionManager(context.CookieStore(r))
		if _, ok := usm.SessionUser(r); !ok {
			http.Error(w, "You must be logged in to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAdmin(r *http.Request) bool {
	usm := session.NewUserSessionManager(context.CookieStore(r))
	user, _ := usm.SessionUser(r)
	return user.IsAdmin
}

func isThreadCreator(r *http.Request) bool {
	tm := models.NewThreadModel(context.DB(r))
	threadID := context.ThreadID(r)
	thread, err := tm.GetThreadByID(threadID)
	if err != nil {
		return false
	}
	usm := session.NewUserSessionManager(context.CookieStore(r))
	user, _ := usm.SessionUser(r)
	return thread.CreatedByEmail == user.Email

}

func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !isAdmin(r) {
			http.Error(w, "You must be an admin to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MustBeAdminOrThreadCreator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !isThreadCreator(r) && !isAdmin(r) {
			http.Error(w, "You must be an admin or creator of the thread to access this link.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
