package model

// User OpenAPI: User
type User struct {
	Id   int    `json:"userId" db:"user_id"`
	Name string `json:"name" db:"name"`
}
