package helpers

import (
	"go-module/models"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var address string = "user=postgres dbname=test password=24052001 sslmode=disable"

func ValidRegisterUsername(s string) bool {
	Re := regexp.MustCompile("^[a-zA-Z0-9]{8,16}$")
	if !Re.MatchString(s) {
		return false
	} else {
		db, err := gorm.Open("postgres", address)
		if err != nil {
			panic("Could not open database to validate register")
		}
		defer db.Close()

		var num int
		Err := db.Table("users").Where("username = ?", s).Count(&num).Error
		if Err != nil {
			panic("Register check crashed")
		}
		return num == 0
	}
}

func ValidRegisterEmail(s string) bool {
	Re := regexp.MustCompile("^[a-zA-Z0-9]+@[a-zA-Z0-9]+.(com|vn)$")
	if !Re.MatchString(s) {
		return false
	} else {
		return true
	}
}

func ValidRegisterPhoneNumber(s string) bool {
	Re := regexp.MustCompile("^[0-9]{10,11}$")
	if !Re.MatchString(s) {
		return false
	} else {
		return true
	}
}

func ValidLogin(x, y string) bool {
	if len(x) == 0 || len(y) == 0 {
		return false
	}

	db, err := gorm.Open("postgres", address)
	if err != nil {
		panic("Could not open database to validate login")
	}
	defer db.Close()

	var user models.User
	Err := db.Table("users").Where("username = ?", x).Find(&user).Error
	if Err != nil {
		return false
	}
	return ComparePasswords(user.Password, y)
}
