package database

import (
	"codigo-fluente/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dns := "root:Lisboa70!@tcp(127.0.0.1:3306)/fluent_admin?charset=utf8mb4&parseTime=True&loc=Local"
	connection, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = connection

	err = connection.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection successfully opened")
}
