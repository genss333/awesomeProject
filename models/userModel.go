package models

type User struct {
	UserId     uint   `gorm:"column:user_id;primaryKey"`
	Username   string `gorm:"column:username"`
	UserEmail  string `gorm:"column:user_email"`
	UserStatus int    `gorm:"column:user_status"`
}
