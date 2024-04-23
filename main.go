package main

import (
	"go-clothes-shop/controllers/authcontroller"
	"go-clothes-shop/controllers/productcontroller"
	"go-clothes-shop/middlewares"
	"go-clothes-shop/models"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	models.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/api/auth/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", authcontroller.Register).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productcontroller.Index).Methods("GET")
	api.HandleFunc("/product/{id}", productcontroller.Show).Methods("GET")
	api.HandleFunc("/auth/logout", authcontroller.Logout).Methods("GET")
	api.Use(middlewares.JWTMiddleware)

	apiAdmin := r.PathPrefix("/api").Subrouter()
	apiAdmin.HandleFunc("/product", productcontroller.Store).Methods("POST")
	apiAdmin.HandleFunc("/product/{id}", productcontroller.Update).Methods("PUT")
	apiAdmin.HandleFunc("/product/{id}", productcontroller.Delete).Methods("DELETE")
	apiAdmin.Use(middlewares.AdminRole)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
