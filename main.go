package main

import (
	"go-clothes-shop/controllers/authcontroller"
	"go-clothes-shop/models"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	models.ConnectDB()

	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", authcontroller.Login).Methods("POST")
	auth.HandleFunc("/register", authcontroller.Register).Methods("POST")

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
