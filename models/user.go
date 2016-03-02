package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{Base{db}}
}

type User struct {
	Username string
}

type UserModel struct {
	Base
}

// GetByUsername returns record by username.
func (u *UserModel) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := u.db.Get(user, "SELECT * FROM users WHERE username=?", username)
	return user, err
}

// Signup creates a new record of user.
func (u *UserModel) Signup(username string) (*User, error) {
	if len(username) == 0 {
		return nil, errors.New("Username cannot be blank.")
	}
	_, err := u.exec("INSERT INTO users(username) VALUES(?)", username)
	if err != nil {
		return nil, err
	}

	return u.GetUserByUsername(username)
}
