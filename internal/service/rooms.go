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

func (s *RoomsService) JoinRoom(userId int) (int, error) {
	return s.repo.JoinRoom(userId)
}

func (s *RoomsService) LeaveRoom(userId int) error {
	return s.repo.LeaveRoom(userId)
}
