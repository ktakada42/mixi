package model

// FriendList OpenAPI: FriendList
type FriendList struct {
	Id   int    `json:"userId" db:"user_id"`
	Name string `json:"name" db:"name"`
}
