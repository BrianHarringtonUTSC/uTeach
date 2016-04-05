// Package models contains models to abstract the db.
package models

import (
	"database/sql/driver"

	"github.com/jmoiron/sqlx"
)

type Base struct {
	db *sqlx.DB
}

func (b *Base) query(tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Rows, error) {
	if tx != nil {
		return tx.Queryx(query, args...)
	}

	return b.db.Queryx(query, args...)
}

func (b *Base) exec(tx *sqlx.Tx, query string, args ...interface{}) (driver.Result, error) {

	if tx != nil {
		return tx.Exec(query, args...)
	}

	return b.db.Exec(query, args...)
}
