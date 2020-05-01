package main

import (
	"fmt"
	"log"
	"net/http"

	"go-module/handlers"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello guys")
	// create our router
	myRouter := mux.NewRouter()

	// our handlers
	myRouter.HandleFunc("/", handlers.HelloWorld).Methods("GET")
	myRouter.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	myRouter.HandleFunc("/users/q/{piece}", handlers.UsersFilter)

	myRouter.HandleFunc("/register", handlers.Register)
	myRouter.HandleFunc("/login", handlers.Login)
	myRouter.HandleFunc("/welcome", handlers.Welcome)
	myRouter.HandleFunc("/update", handlers.GetUpdate).Methods("GET")
	myRouter.HandleFunc("/update", handlers.PostUpdate).Methods("POST")
	myRouter.HandleFunc("/logout", handlers.Logout).Methods("GET")
	myRouter.HandleFunc("/release", handlers.Release)

	// connect to the url
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
