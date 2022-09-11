package service

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"

	"problem1/model"
	"problem1/repository"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListService interface {
	CheckUserExist(userId int) (bool, error)
	InsertUserLink(ulfr *model.UserLinkForRequest) error
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

func (s *friendListService) CheckUserExist(userId int) (bool, error) {
	return s.flr.CheckUserExist(userId)
}

func (s *friendListService) InsertUserLink(ulfr *model.UserLinkForRequest) error {
	if err := s.flr.CheckUserLink(ulfr.User1Id, ulfr.User2Id, ulfr.Table); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.flr.InsertUserLink(ulfr.User1Id, ulfr.User2Id, ulfr.Table)
		}

		return err
	}

	return nil
}

func (s *friendListService) GetFriendListByUserId(c echo.Context) (*model.FriendList, error) {
	userId := c.Get("userId").(int)

	blockUsers, err := s.flr.GetBlockUsersIdList(userId)
	if err != nil {
		return nil, err
	}
	if len(blockUsers) == 0 {
		return s.flr.GetFriendListByUserId(userId)
	}

	return s.flr.GetFriendListByUserIdExcludingBlockUsers(userId, blockUsers)
}

func (s *friendListService) GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error) {
	userId := c.Get("userId").(int)

	oneHopFriends, err := s.flr.GetOneHopFriendsUserIdList(userId)
	if err != nil {
		return nil, err
	}
	if len(oneHopFriends) == 0 {
		return &model.FriendList{Friends: nil}, nil
	}

	blockUsers, err := s.flr.GetBlockUsersIdList(userId)
	if err != nil {
		return nil, err
	}

	excludeUsers := append(oneHopFriends, blockUsers...)

	return s.flr.GetFriendListOfFriendsByUserId(userId, excludeUsers)
}

func (s *friendListService) GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) (*model.FriendList, error) {
	userId := c.Get("userId").(int)
	limit := c.Get("limit").(int)
	offset := c.Get("offset").(int)

	oneHopFriends, err := s.flr.GetOneHopFriendsUserIdList(userId)
	if err != nil {
		return nil, err
	}
	if len(oneHopFriends) == 0 {
		return &model.FriendList{Friends: nil}, nil
	}

	blockUsers, err := s.flr.GetBlockUsersIdList(userId)
	if err != nil {
		return nil, err
	}

	excludeUsers := append(oneHopFriends, blockUsers...)

	return s.flr.GetFriendListOfFriendsByUserIdWithPaging(userId, excludeUsers, limit, offset)
}
