package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	return "/users/" + u.Email
}

// UserModel handles getting and creating users.
type UserModel struct {
	Base
}

// NewUserModel returns a new user model.
func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{Base{db}}
}

var usersBuilder = squirrel.Select("* FROM users")

// Find gets all users filtered by wheres.
func (um *UserModel) Find(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) ([]*User, error) {
	selectBuilder := um.addWheresToBuilder(usersBuilder, wheres...)
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	var users []*User
	err = um.sel(tx, &users, query, args...)
	return users, errors.Wrap(err, "select error")
}

// FindOne gets the user filtered by wheres.
func (um *UserModel) FindOne(tx *sqlx.Tx, wheres ...squirrel.Sqlizer) (*User, error) {
	users, err := um.Find(tx, wheres...)
	if err != nil {
		return nil, err
	}

	switch len(users) {
	case 0:
		return nil, sql.ErrNoRows
	case 1:
		return users[0], nil
	default:
		msg := fmt.Sprintf("expected 1, got %d", len(users))
		return nil, errors.New(msg)
	}
}

// AddUser adds a new user.
func (um *UserModel) AddUser(tx *sqlx.Tx, email, name string) (*User, error) {
	if email == "" || name == "" {
		return nil, InputError{"email and/or name cannot be empty"}
	}

	email = strings.ToLower(email)
	name = strings.Title(name)
	result, err := um.exec(tx, "INSERT INTO users(email, name) VALUES(?, ?)", email, name)
	if err != nil {
		return nil, errors.Wrap(err, "exec error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "last inserted id error")
	}

	user, err := um.FindOne(tx, squirrel.Eq{"users.id": id})
	return user, errors.Wrap(err, "find one error")
}
