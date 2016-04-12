package models

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// User represents a user in the app.
type User struct {
	ID      int64
	Email   string
	Name    string
	IsAdmin bool `db:"is_admin"`
}

// URL returns the unique URL for a user.
func (u *User) URL() string {
	return "/user/" + u.Email
}

// UserModel handles getting and creating users.
type UserModel struct {
	Base
}

// NewUserModel returns a new user model.
func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{Base{db}}
}

// GetUserByID gets a user by the id.
func (um *UserModel) GetUserByID(tx *sqlx.Tx, id int64) (*User, error) {
	user := new(User)
	err := um.Get(tx, user, "SELECT * FROM users WHERE id=?", id)
	return user, err
}

// GetUserByEmail gets a user by email.
func (um *UserModel) GetUserByEmail(tx *sqlx.Tx, email string) (*User, error) {
	email = strings.ToLower(email)

	user := new(User)
	err := um.Get(tx, user, "SELECT * FROM users WHERE email=?", email)
	return user, err
}

// AddUser adds a new user.
func (um *UserModel) AddUser(tx *sqlx.Tx, email, name string) (*User, error) {
	if email == "" || name == "" {
		return nil, InputError{"email and/or name cannot be empty"}
	}

	email = strings.ToLower(email)
	name = strings.Title(name)
	_, err := um.Exec(tx, "INSERT INTO users(email, name) VALUES(?, ?)", email, name)
	if err != nil {
		return nil, err
	}

	return um.GetUserByEmail(tx, email)
}
