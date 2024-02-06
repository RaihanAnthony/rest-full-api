package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:0987@tcp(localhost:3306)/rest_full_api"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	DB = db
}