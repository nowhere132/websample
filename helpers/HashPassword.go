package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPass, plainPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
