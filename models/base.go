// Package models provides models for the uTeach app.
package models

import (
	"database/sql/driver"
	"github.com/jmoiron/sqlx"
)

type Base struct {
	db *sqlx.DB
}

func (b *Base) exec(query string, params ...interface{}) (driver.Result, error) {
	stmt, err := b.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(params...)
}
