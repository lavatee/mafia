package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lavatee/mafia"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersTable       = "users"
	roomsTable       = "rooms"
	superpowersTable = "superpowers"
	playersTable     = "players"
	requestsTable    = "requests"
)

type Auth interface {
	SignUp(email string, name string, passwordHash string) (int, error)
	SignIn(email string, passwordHash string) (mafia.User, error)
}
type Rooms interface {
	JoinRoom(userId int) (int, error)
	LeaveRoom(userId int) error
}
type Friends interface {
	GetFriends(id int) ([]MongoFriend, error)
	AddFriend(userId int, friendId int) error
	DeleteFriend(userId int, friendId int) error
}
type FriendshipRequests interface {
}
type PlayersInRoom interface {
}
type Superpowers interface {
	NewSuperpower(userId int, name string) (int, error)
}
type Repository struct {
	Auth
	Rooms
	Friends
	FriendshipRequests
	PlayersInRoom
	Superpowers
}

func NewRepository(db *sqlx.DB, mongo *mongo.Collection, rdb *redis.Client) *Repository {
	return &Repository{
		Auth:               NewAuthPostgres(db, mongo),
		Rooms:              NewRoomsPostgres(db),
		Friends:            NewFriendsMongo(rdb, mongo, db),
		FriendshipRequests: NewRequestsPostgres(db),
		PlayersInRoom:      NewPlayersPostgres(db),
		Superpowers:        NewSuperpowersPostgres(db),
	}
}
