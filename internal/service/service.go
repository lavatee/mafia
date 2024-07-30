package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/mafia/internal/repository"
)

type Auth interface {
	HashPassword(password string) string
	SignUp(email string, name string, password string) (int, error)
	SignIn(email string, password string) (string, string, error)
	NewToken(claims jwt.MapClaims) (string, error)
	Refresh(token string) (string, string, error)
}
type Rooms interface {
}
type Friends interface {
	GetFriends(id int) ([]repository.MongoFriend, error)
	AddFriend(userId int, friendId int) error
}
type FriendshipRequests interface {
}
type PlayersInRoom interface {
}
type Superpowers interface {
}
type Service struct {
	Auth
	Rooms
	Friends
	FriendshipRequests
	PlayersInRoom
	Superpowers
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:               NewAuthService(repo),
		Rooms:              NewRoomsService(repo),
		Friends:            NewFriendsService(repo),
		FriendshipRequests: NewRequestsService(repo),
		PlayersInRoom:      NewPlayersService(repo),
		Superpowers:        NewSuperpowersService(repo),
	}
}
