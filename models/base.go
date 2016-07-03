// Package models contains models to manage communications with the db.
// All models methods should have a tx as the first parameter. This allows handlers that use multiple calls / tables to
// use a single transaction and then commit that so that all actions are committed or none are.
// If a handler is using a single action, they can simply pass in nil for the tx.
package models

import (
	"database/sql/driver"
	"regexp"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var singleWordAlphaNumRegex = regexp.MustCompile(`^[[:alnum:]]+(_[[:alnum:]]+)*$`)

// InputError is an error returned when the user supplied input was not valid.
type InputError struct {
	Message string
}

// Error returns the input error's message.
func (ie InputError) Error() string {
	return ie.Message
}

// Base is the base model for all other models to embed.
// It has common helpers and functionality that all models can use.
type Base struct {
	db *sqlx.DB
}

func (b *Base) exec(tx *sqlx.Tx, query string, args ...interface{}) (driver.Result, error) {
	if tx != nil {
		return tx.Exec(query, args...)
	}

	return b.db.Exec(query, args...)
}

func (b *Base) get(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	if tx != nil {
		return tx.Get(dest, query, args...)
	}

	return b.db.Get(dest, query, args...)
}

func (b *Base) sel(tx *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	if tx != nil {
		return tx.Select(dest, query, args...)
	}

	return b.db.Select(dest, query, args...)
}

func (b *Base) query(tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Rows, error) {
	if tx != nil {
		return tx.Queryx(query, args...)
	}

	return b.db.Queryx(query, args...)
}

func (b *Base) addWheresToBuilder(selectBuilder squirrel.SelectBuilder, wheres ...squirrel.Sqlizer) squirrel.SelectBuilder {
	for _, where := range wheres {
		selectBuilder = selectBuilder.Where(where)
	}

	return selectBuilder
}

func (b *Base) queryWhere(tx *sqlx.Tx, selectBuilder squirrel.SelectBuilder, wheres ...squirrel.Sqlizer) (*sqlx.Rows, error) {
	selectBuilder = b.addWheresToBuilder(selectBuilder, wheres...)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	return b.query(tx, query, args...)
}
