package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"primary_key"`
	Email       string
	PhoneNumber string
	Password    string
}
