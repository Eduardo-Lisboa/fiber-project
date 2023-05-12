package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty" gorm:"unique"`
	Password  []byte
}
