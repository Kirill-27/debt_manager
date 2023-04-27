package service

import "github.com/kirill-27/debt_manager/pkg/repository"

type FriendsService struct {
	repo repository.Friends
}

func NewFriendsService(repo repository.Friends) *FriendsService {
	return &FriendsService{repo: repo}
}

func (s *FriendsService) AddFriend(myId int, friendId int) error {
	return s.repo.AddFriend(myId, friendId)
}
