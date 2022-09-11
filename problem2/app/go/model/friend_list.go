package model

// UserLinkForRequest OpenAPI: UserLinkForRequest
type UserLinkForRequest struct {
	User1Id int    `json:"user1Id" db:"user1_id"`
	User2Id int    `json:"user2Id" db:"user2_id"`
	Table   string `json:"table"`
}

// Friend OpenAPI: Friend
type Friend struct {
	UserId int    `json:"userId" db:"user_id"`
	Name   string `json:"name" db:"name"`
}

// FriendList OpenAPI: FriendList
type FriendList struct {
	Friends []*Friend `json:"friends"`
}
