package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/mafia"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthPostgres struct {
	db  *sqlx.DB
	mdb *mongo.Collection
}

type MongoFriend struct {
	Id   int
	Name string
}

type MongoUser struct {
	Id      int           `json:"Id"`
	Friends []MongoFriend `json:"Friends"`
}

func NewAuthPostgres(db *sqlx.DB, mdb *mongo.Collection) *AuthPostgres {
	return &AuthPostgres{
		db:  db,
		mdb: mdb,
	}
}

func (r *AuthPostgres) SignUp(email string, name string, passwordHash string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, name, password_hash, coins) values ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, email, name, passwordHash, 0)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	newMongoUser := MongoUser{
		Id:      id,
		Friends: []MongoFriend{},
	}
	_, err := r.mdb.InsertOne(context.Background(), newMongoUser)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) SignIn(email string, passwordHash string) (mafia.User, error) {
	var user mafia.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, passwordHash)
	return user, err
}
