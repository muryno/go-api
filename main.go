package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-api/app"
	"go-api/controller"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", controller.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controller.UserLogin).Methods("POST")
	router.HandleFunc("/api/user/{user_id}", controller.GetUser).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}

	log.Fatal(err)

}
