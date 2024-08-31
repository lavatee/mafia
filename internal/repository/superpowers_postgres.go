package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SuperpowersPostgres struct {
	db *sqlx.DB
}

func NewSuperpowersPostgres(db *sqlx.DB) *SuperpowersPostgres {
	return &SuperpowersPostgres{
		db: db,
	}
}

func (r *SuperpowersPostgres) NewSuperpower(userId int, name string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, name) values ($1, $2) RETURNING id", superpowersTable)
	row := r.db.QueryRow(query, userId, name)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
