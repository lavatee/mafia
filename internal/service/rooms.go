package service

import "github.com/lavatee/mafia/internal/repository"

type RoomsService struct {
	repo *repository.Repository
}

func NewRoomsService(repo *repository.Repository) *RoomsService {
	return &RoomsService{
		repo: repo,
	}
}
