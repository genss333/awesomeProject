package database

import (
	"awesomeProject/utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var username = utils.GoDotEnvVariable("DB_USERNAME")
var host = utils.GoDotEnvVariable("DB_HOST")
var port = utils.GoDotEnvVariable("DB_PORT")
var database = utils.GoDotEnvVariable("DB_DATABASE")

var _ *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s", username, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
