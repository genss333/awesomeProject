package models

type User struct {
	UserId     int       `gorm:"column:user_id;primaryKey;" json:"user_id"`
	Username   string    `gorm:"column:username" json:"username"`
	Password   string    `gorm:"column:password" json:"password"`
	UserEmail  string    `gorm:"column:user_email" json:"user_email"`
	UserStatus int       `gorm:"column:user_status;default:1" json:"user_status"`
	Books      Book      `gorm:"foreignKey:UserID" json:"books"`
	UserImages UserImage `gorm:"foreignKey:UserID" json:"user_images"`
}
