package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"problem1/model"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListRepository interface {
	CheckUserExist(userId int) (bool, error)
	GetOneHopFriendsUserIdList(userId int) ([]int, error)
	GetBlockUsersIdList(userId int) ([]int, error)
	GetFriendListByUserId(userId int) (*model.FriendList, error)
	GetFriendListByUserIdExcludingBlockUsers(userId int, blockUsers []int) (*model.FriendList, error)
	GetFriendListOfFriendsByUserId(userId int, excludeUsers []int) (*model.FriendList, error)
	GetFriendListOfFriendsByUserIdWithPaging(userId int, excludeUsers []int, limit, offset int) (*model.FriendList, error)
}

type friendListRepository struct {
	db *sql.DB
}

func NewFriendListRepository(db *sql.DB) FriendListRepository {
	return &friendListRepository{
		db: db,
	}
}

func (r *friendListRepository) CheckUserExist(userId int) (bool, error) {
	const q = `
	SELECT user_id, name
	FROM users
	WHERE user_id = ?`

	row := r.db.QueryRow(q, userId)

	user := &model.Friend{}
	if err := row.Scan(&user.UserId, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *friendListRepository) GetOneHopFriendsUserIdList(userId int) ([]int, error) {
	const q = `
	SELECT user2_id
	FROM friend_link
	WHERE user1_id = ?;`

	rows, err := r.db.Query(q, userId)
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

	return oneHopFriends, nil
}

func (r *friendListRepository) GetBlockUsersIdList(userId int) ([]int, error) {
	const q = `
	SELECT user2_id
	FROM block_list
	WHERE user1_id = ?;`

	rows, err := r.db.Query(q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		blockUsers []int
		blockUser  int
	)
	for rows.Next() {
		if err := rows.Scan(&blockUser); err != nil {
			return nil, err
		}

		blockUsers = append(blockUsers, blockUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blockUsers, nil
}

func (r *friendListRepository) GetFriendListByUserId(userId int) (*model.FriendList, error) {
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

	var friends []*model.Friend
	for rows.Next() {
		friend := &model.Friend{}
		if err := rows.Scan(&friend.UserId, &friend.Name); err != nil {
			return nil, err
		}

		friends = append(friends, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &model.FriendList{Friends: friends}, nil
}

func (r *friendListRepository) GetFriendListByUserIdExcludingBlockUsers(userId int, blockUsers []int) (*model.FriendList, error) {
	const q = `
	SELECT U.user_id, U.name
	FROM users AS U INNER JOIN friend_link AS FL
	ON U.user_id = FL.user2_id
	WHERE FL.user1_id = ?
	AND U.user_id NOT IN (?);`

	dbx := sqlx.NewDb(r.db, "mysql")

	query, args, err := sqlx.In(q, userId, blockUsers)
	if err != nil {
		return nil, err
	}

	var friends []*model.Friend
	if err := dbx.Select(&friends, query, args...); err != nil {
		return nil, err
	}

	return &model.FriendList{Friends: friends}, nil
}

func (r *friendListRepository) GetFriendListOfFriendsByUserId(userId int, excludeUsers []int) (*model.FriendList, error) {
	const q = `
	SELECT DISTINCT U.user_id, U.name
	FROM users AS U
	INNER JOIN friend_link AS FL
	ON U.user_id = FL.user2_id
	INNER JOIN friend_link AS FL2
	ON FL.user1_id = FL2.user2_id
	WHERE FL2.user1_id = ?
	AND U.user_id NOT IN (?);`

	dbx := sqlx.NewDb(r.db, "mysql")

	query, args, err := sqlx.In(q, userId, excludeUsers)
	if err != nil {
		return nil, err
	}

	var friends []*model.Friend
	if err := dbx.Select(&friends, query, args...); err != nil {
		return nil, err
	}

	return &model.FriendList{Friends: friends}, nil
}

func (r *friendListRepository) GetFriendListOfFriendsByUserIdWithPaging(userId int, excludeUsers []int, limit, offset int) (*model.FriendList, error) {
	const q = `
	SELECT DISTINCT U.user_id, U.name
	FROM users AS U
	INNER JOIN friend_link AS FL
	ON U.user_id = FL.user2_id
	INNER JOIN friend_link AS FL2
	ON FL.user1_id = FL2.user2_id
	WHERE FL2.user1_id = ?
	AND U.user_id NOT IN (?)
	LIMIT ? OFFSET ?;`

	dbx := sqlx.NewDb(r.db, "mysql")

	query, args, err := sqlx.In(q, userId, excludeUsers, limit, offset)
	if err != nil {
		return nil, err
	}

	var friends []*model.Friend
	if err := dbx.Select(&friends, query, args...); err != nil {
		return nil, err
	}

	return &model.FriendList{Friends: friends}, nil
}
