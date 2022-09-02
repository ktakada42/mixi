package repository

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"problem1/model"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListRepository interface {
	GetFriendListByUserId(c echo.Context) ([]*model.FriendList, error)
}

type friendListRepository struct {
	db *sql.DB
}

func NewFriendListRepository(db *sql.DB) FriendListRepository {
	return &friendListRepository{
		db: db,
	}
}

func (r *friendListRepository) GetFriendListByUserId(c echo.Context) ([]*model.FriendList, error) {
	userId := c.QueryParam("userId")

	const q = `
SELECT U.user_id, U.name
FROM users AS U INNER JOIN friend_link AS FL
ON U.user_id = FL.user2_id
WHERE FL.user1_id = ?`

	rows, err := r.db.Query(q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendLists []*model.FriendList
	for rows.Next() {
		friendList := &model.FriendList{}
		if err := rows.Scan(&friendList.Id, &friendList.Name); err != nil {
			return nil, err
		}

		friendLists = append(friendLists, friendList)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friendLists, nil
}
