package repository

import "github.com/jmoiron/sqlx"

type PlayersPostgres struct {
	db *sqlx.DB
}

func NewPlayersPostgres(db *sqlx.DB) *PlayersPostgres {
	return &PlayersPostgres{
		db: db,
	}
}
