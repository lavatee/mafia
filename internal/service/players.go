package service

import "github.com/lavatee/mafia/internal/repository"

type PlayersService struct {
	repo *repository.Repository
}

func NewPlayersService(repo *repository.Repository) *PlayersService {
	return &PlayersService{
		repo: repo,
	}
}
