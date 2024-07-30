package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(d PostgresDB) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s username=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.Username, d.Password, d.DBName, d.SSLMode))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
