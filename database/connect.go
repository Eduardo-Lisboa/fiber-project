package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() {
	dns := "root:Lisboa70!@tcp(127.0.0.1:3306)/fluent_admin?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection successfully opened")
}
