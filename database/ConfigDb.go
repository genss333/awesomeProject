package database

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var username = utils.GoDotEnvVariable("DB_USERNAME")
var password = utils.GoDotEnvVariable("DB_PASSWORD")
var host = utils.GoDotEnvVariable("DB_HOST")
var port = utils.GoDotEnvVariable("DB_PORT")
var database = utils.GoDotEnvVariable("DB_DATABASE")

var _ *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables() {
	db, err := Connect()
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(TableList()...)
	if err != nil {
		return
	}
}

func TableList() []interface{} {
	model := []interface{}{&models.User{}, &models.Book{}, &models.UserImage{}}
	return model
}
