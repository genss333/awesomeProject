package models

type User struct {
	UserId     int    `json:"user_id"`
	Username   string `json:"username"`
	UserEmail  string `json:"user_email"`
	UserStatus int    `json:"user_status"`
}
