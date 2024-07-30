package service

import "github.com/lavatee/mafia/internal/repository"

type RequestsService struct {
	repo *repository.Repository
}

func NewRequestsService(repo *repository.Repository) *RequestsService {
	return &RequestsService{
		repo: repo,
	}
}
