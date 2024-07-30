package repository

import "github.com/jmoiron/sqlx"

type SuperpowersPostgres struct {
	db *sqlx.DB
}

func NewSuperpowersPostgres(db *sqlx.DB) *SuperpowersPostgres {
	return &SuperpowersPostgres{
		db: db,
	}
}
