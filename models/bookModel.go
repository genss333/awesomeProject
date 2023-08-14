package models

type Book struct {
	BookId  int    `gorm:"column:book_id;primary_key"`
	UserId  uint   `gorm:"column:user_id"`
	Address string `gorm:"column:address"`
	Tel     string `gorm:"column:tel"`
	PId     string `gorm:"column:p_id"`
}
