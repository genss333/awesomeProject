package models

type Book struct {
	BookId  int    `gorm:"column:book_id;primary_key" json:"book_id"`
	Address string `gorm:"column:address" json:"address"`
	Tel     string `gorm:"column:tel" json:"tel"`
	PId     string `gorm:"column:p_id" json:"p_id"`
	UserID  int    `gorm:"column:user_id"`
}
