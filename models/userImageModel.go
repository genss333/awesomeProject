package models

type UserImage struct {
	UserImageId int    `gorm:"column:user_image_id;primaryKey" json:"user_image_id"`
	Image       string `gorm:"column:image" json:"image"`
	UserID      int    `gorm:"column:user_id"`
}
