package repository

import "github.com/jmoiron/sqlx"

type RoomsPostgres struct {
	db *sqlx.DB
}

func NewRoomsPostgres(db *sqlx.DB) *RoomsPostgres {
	return &RoomsPostgres{
		db: db,
	}
}
