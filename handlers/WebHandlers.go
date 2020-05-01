package handlers

import (
	"encoding/json"
	"fmt"
	"go-module/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var address string = "user=postgres dbname=test password=24052001 sslmode=disable"

func init() {
	db, err := gorm.Open("postgres", address)
	if err != nil {
		panic("Could not init database")
	}
	defer db.Close()

	db.AutoMigrate(&models.User{})
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This link is to test my get method")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", address)
	if err != nil {
		panic("Could not open database to show users")
	}
	defer db.Close()

	var users []models.User
	db.Table("users").Find(&users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UsersFilter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	piece := vars["piece"]
	// fmt.Println(piece)

	db, err := gorm.Open("postgres", address)
	if err != nil {
		panic("Could not open database for filter")
	}
	defer db.Close()

	var users []models.User
	db.Table("users").Where("username LIKE ?", "%"+piece+"%").Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
