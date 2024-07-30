package service

import "github.com/lavatee/mafia/internal/repository"

type FriendsService struct {
	repo *repository.Repository
}

func NewFriendsService(repo *repository.Repository) *FriendsService {
	return &FriendsService{
		repo: repo,
	}
}
func (s *FriendsService) GetFriends(id int) ([]repository.MongoFriend, error) {
	return s.repo.GetFriends(id)
}

func (s *FriendsService) AddFriend(userId int, friendId int) error {
	return s.repo.AddFriend(userId, friendId)
}
