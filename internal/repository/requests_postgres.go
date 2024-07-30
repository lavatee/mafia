package repository

import "github.com/jmoiron/sqlx"

type RequestsPostgres struct {
	db *sqlx.DB
}

func NewRequestsPostgres(db *sqlx.DB) *RequestsPostgres {
	return &RequestsPostgres{
		db: db,
	}
}
