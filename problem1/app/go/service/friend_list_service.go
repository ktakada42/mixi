package service

import (
	"github.com/labstack/echo/v4"

	"problem1/model"
	"problem1/repository"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListService interface {
	GetFriendListByUserId(c echo.Context) ([]*model.FriendList, error)
}

type friendListService struct {
	flr repository.FriendListRepository
}

func NewFriendListService(flr repository.FriendListRepository) FriendListService {
	return &friendListService{
		flr: flr,
	}
}

func (s *friendListService) GetFriendListByUserId(c echo.Context) ([]*model.FriendList, error) {
	return s.flr.GetFriendListByUserId(c)
}
