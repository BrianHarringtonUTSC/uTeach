// Package models contains models to manage communicating with the db.
// All models methods have a tx as the first parameter. This allows handlers that use multiple calls / tables to
// use a single transaction and then commit that so that all actions are committed or none are.
// If a handler is using a single action, they can simply pass in nil for the tx.
package models

import (
	"database/sql/driver"
	"regexp"

	"github.com/jmoiron/sqlx"
)

var singleWordAlphaNumRegex = regexp.MustCompile(`^[[:alnum:]]+(_[[:alnum:]]+)*$`)

// Base is the base model for all other models to embed.
// It has common helpers and functionality that all models can use.
type Base struct {
	db *sqlx.DB
}

// InputError is an error returned when the user supplied input was not valid.
type InputError struct {
	Message string
}

// Error returns the input error's message.
func (ie InputError) Error() string {
	return ie.Message
}

// Exec execs with the tx if not null else with the db.
func (b *Base) Exec(tx *sqlx.Tx, query string, args ...interface{}) (driver.Result, error) {

	if tx != nil {
		return tx.Exec(query, args...)
	}

	return b.db.Exec(query, args...)
}

// Get execs with the tx if not null else with the db.
func (b *Base) Get(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	if tx != nil {
		return tx.Get(dest, query, args...)
	}

	return b.db.Get(dest, query, args...)
}

// Select execs with the tx if not null else with the db.
func (b *Base) Select(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	if tx != nil {
		return tx.Select(dest, query, args...)
	}

	return b.db.Select(dest, query, args...)
}

// Query execs with the tx if not null else with the db.
func (b *Base) Query(tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Rows, error) {
	if tx != nil {
		return tx.Queryx(query, args...)
	}

	return b.db.Queryx(query, args...)
}
