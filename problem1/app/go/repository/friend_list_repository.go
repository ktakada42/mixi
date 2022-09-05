package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"problem1/model"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListRepository interface {
	CheckUserExist(c echo.Context) (bool, error)
	GetFriendListByUserId(c echo.Context) (*model.FriendList, error)
	GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error)
}

type friendListRepository struct {
	db *sql.DB
}

func NewFriendListRepository(db *sql.DB) FriendListRepository {
	return &friendListRepository{
		db: db,
	}
}

func (r *friendListRepository) CheckUserExist(c echo.Context) (bool, error) {
	userId := c.QueryParam("userId")

	const q = `
SELECT user_id, name
FROM users
WHERE user_id = ?`

	row := r.db.QueryRow(q, userId)

	user := &model.User{}
	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *friendListRepository) GetFriendListByUserId(c echo.Context) (*model.FriendList, error) {
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

	var friends []*model.User
	for rows.Next() {
		friend := &model.User{}
		if err := rows.Scan(&friend.Id, &friend.Name); err != nil {
			return nil, err
		}

		friends = append(friends, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &model.FriendList{Friends: friends}, nil
}

func (r *friendListRepository) GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error) {
	userId := c.QueryParam("userId")

	// get 1hop friends
	const q1 = `
SELECT user2_id
FROM friend_link
WHERE user1_id = ?;`

	rows, err := r.db.Query(q1, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		oneHopFriends []int
		oneHopFriend  int
	)
	for rows.Next() {
		if err := rows.Scan(&oneHopFriend); err != nil {
			return nil, err
		}

		oneHopFriends = append(oneHopFriends, oneHopFriend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get 2hop Friends except for 1hop friends
	const q2 = `
	SELECT U.user_id, U.name
		FROM users AS U
	INNER JOIN friend_link AS FL
	ON U.user_id = FL.user2_id
	INNER JOIN friend_link AS FL2
	ON FL.user1_id = FL2.user2_id
	WHERE FL2.user1_id = ?
	AND U.user_id NOT IN (?);`

	dbx := sqlx.NewDb(r.db, "mysql")

	q, args, err := sqlx.In(q2, userId, oneHopFriends)
	if err != nil {
		return nil, err
	}

	var friends []*model.User
	if err := dbx.Select(&friends, q, args...); err != nil {
		return nil, err
	}

	return &model.FriendList{Friends: friends}, nil
}
