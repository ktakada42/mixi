package service

import (
	"github.com/labstack/echo/v4"

	"problem1/model"
	"problem1/repository"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListService interface {
	CheckUserExist(c echo.Context) (bool, error)
	GetFriendListByUserId(c echo.Context) (*model.FriendList, error)
	GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error)
	GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) (*model.FriendList, error)
}

type friendListService struct {
	flr repository.FriendListRepository
}

func NewFriendListService(flr repository.FriendListRepository) FriendListService {
	return &friendListService{
		flr: flr,
	}
}

func (s *friendListService) CheckUserExist(c echo.Context) (bool, error) {
	return s.flr.CheckUserExist(c)
}

func (s *friendListService) GetFriendListByUserId(c echo.Context) (*model.FriendList, error) {
	return s.flr.GetFriendListByUserId(c)
}

func (s *friendListService) GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error) {
	return s.flr.GetFriendListOfFriendsByUserId(c)
}
func (s *friendListService) GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) (*model.FriendList, error) {
	return s.flr.GetFriendListOfFriendsByUserIdWithPaging(c)
}
