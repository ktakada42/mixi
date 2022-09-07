package model

// Friend OpenAPI: Friend
type Friend struct {
	UserId int    `json:"userId" db:"user_id"`
	Name   string `json:"name" db:"name"`
}

// FriendList OpenAPI: FriendList
type FriendList struct {
	Friends []*Friend `json:"friends"`
}
