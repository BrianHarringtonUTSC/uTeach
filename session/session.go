// Package session provides functionality for user sessions and cookie management.
package session

import (
	"database/sql"
	"encoding/gob"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"net/http"

	"github.com/umairidris/uTeach/models"
)

const (
	userSessionName = "user-session"
	userKey         = "user"
)

func init() {
	// allows user to be encoded so that it can be stored in a session
	gob.Register(&models.User{})
}

type UserSession struct {
	store *sessions.CookieStore
}

func NewUserSessionManager(store *sessions.CookieStore) *UserSession {
	return &UserSession{store}
}

// getUserSession gets the session containing the user.
func (us *UserSession) get(r *http.Request) (*sessions.Session, error) {
	return us.store.Get(r, userSessionName)
}

// New creates a new session and stores the User containing the
func (us *UserSession) New(w http.ResponseWriter, r *http.Request, email string, name string, db *sqlx.DB) error {
	session, err := us.get(r)
	if err != nil {
		return err
	}

	u := models.NewUserModel(db)
	user, err := u.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		user, err = u.Signup(email, name)
	}
	if err != nil {
		return err
	}

	session.Values[userKey] = user
	session.Options.HttpOnly = true
	return session.Save(r, w)
}

// User gets the user stored in the user session.
// If there is a User stored in the session and can be retrieved it returns the user and true, else the boolean will be
// false.
func (us *UserSession) SessionUser(r *http.Request) (*models.User, bool) {
	session, err := us.get(r)
	if err != nil {
		return nil, false
	}

	um, ok := session.Values[userKey]
	if !ok {
		return nil, false
	}

	user, ok := um.(*models.User)
	return user, ok
}

// Delete deletes the user session.
func (us *UserSession) Delete(w http.ResponseWriter, r *http.Request) error {
	session, err := us.get(r)
	if err != nil {
		return err
	}
	delete(session.Values, userKey)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
