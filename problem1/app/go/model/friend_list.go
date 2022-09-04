package model

// User OpenAPI: User
type User struct {
	Id   int    `json:"userId" db:"user_id"`
	Name string `json:"name" db:"name"`
}

// FriendList OpenAPI: FriendList
type FriendList struct {
	Friends []*User `json:"friends"`
}
