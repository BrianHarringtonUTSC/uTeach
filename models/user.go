package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{Base{db}}
}

type User struct {
	ID      int64
	Email   string
	Name    string
	IsAdmin bool `db:"is_admin"`
}

type UserModel struct {
	Base
}

// GetUserByEmail returns record by email.
func (um *UserModel) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT * FROM users WHERE email=?"
	err := um.db.Get(user, query, email)
	return user, err
}

// Signup creates a new record of user.
func (um *UserModel) Signup(email, name string) (*User, error) {
	if email == "" || name == "" {
		return nil, errors.New(".")
	}

	email = strings.ToLower(email)
	name = strings.Title(name)
	_, err := um.exec("INSERT INTO users(email, name) VALUES(?, ?)", email, name)
	if err != nil {
		return nil, err
	}

	return um.GetUserByEmail(email)
}
