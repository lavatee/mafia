package service

import "github.com/lavatee/mafia/internal/repository"

type SuperpowersService struct {
	repo *repository.Repository
}

func NewSuperpowersService(repo *repository.Repository) *SuperpowersService {
	return &SuperpowersService{
		repo: repo,
	}
}

func (s *SuperpowersService) NewSuperpower(userId int, name string) (int, error) {
	return s.repo.NewSuperpower(userId, name)
}
