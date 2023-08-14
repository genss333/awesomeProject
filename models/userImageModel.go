package models

type UserImage struct {
	UserImageId int    `gorm:"column:user_image_id;primaryKey"`
	UserId      uint   `gorm:"column:user_id"`
	Image       string `gorm:"column:image"`
}
