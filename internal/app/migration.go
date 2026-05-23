package app

import (
	"github.com/jmoiron/sqlx"

	"github.com/pressly/goose/v3"
)

func MigrateUp(db *sqlx.DB, dir string) error {

	if err := goose.Up(db.DB, dir); err != nil {
		return err
	}
	return nil
}

func MigrateDown(db *sqlx.DB, dir string) error {
	if err := goose.Down(db.DB, dir); err != nil {
		return err
	}
	return nil
}
