package usecase

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"problem1/model"
	"problem1/service"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListUseCase interface {
	GetFriendListByUserId(c echo.Context) ([]*model.User, error)
}

type friendListUseCase struct {
	db  *sql.DB
	fls service.FriendListService
}

func NewFriendListUseCase(db *sql.DB, fls service.FriendListService) FriendListUseCase {
	return &friendListUseCase{
		db:  db,
		fls: fls,
	}
}

func (u *friendListUseCase) GetFriendListByUserId(c echo.Context) ([]*model.User, error) {
	return u.fls.GetFriendListByUserId(c)
}
