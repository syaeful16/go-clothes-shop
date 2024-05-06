package main

import (
	"go-clothes-shop/controllers/authcontroller"
	"go-clothes-shop/controllers/cartcontroller"
	"go-clothes-shop/controllers/datausercontroller"
	"go-clothes-shop/controllers/detailproductcontroller"
	"go-clothes-shop/controllers/imagecontroller"
	"go-clothes-shop/controllers/productcontroller"
	"go-clothes-shop/controllers/useraddresscontroller"
	"go-clothes-shop/controllers/userdetailcontroller"
	"go-clothes-shop/middlewares"
	"go-clothes-shop/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	models.ConnectDB()

	r := mux.NewRouter()
	r.HandleFunc("/api/auth/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", authcontroller.Register).Methods("POST")

	// get all product
	r.HandleFunc("/api/products", productcontroller.Index).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	// check user
	api.HandleFunc("/user", authcontroller.User).Methods("GET")

	api.HandleFunc("/auth/logout", authcontroller.Logout).Methods("GET")

	// user address
	api.HandleFunc("/user-address", useraddresscontroller.Index).Methods("GET")
	api.HandleFunc("/user-address/{id}", useraddresscontroller.Show).Methods("GET")
	api.HandleFunc("/user-address", useraddresscontroller.Store).Methods("POST")
	api.HandleFunc("/user-address/{id}", useraddresscontroller.Update).Methods("PUT")
	api.HandleFunc("/user-address/{id}", useraddresscontroller.Delete).Methods("DELETE")

	// user detail
	api.HandleFunc("/user-detail", userdetailcontroller.Index).Methods("GET")
	api.HandleFunc("/user-detail", userdetailcontroller.Store).Methods("POST")
	api.HandleFunc("/user-detail", userdetailcontroller.Update).Methods("PUT")
	api.HandleFunc("/user-detail/{id}", userdetailcontroller.Delete).Methods("DELETE")

	// product
	api.HandleFunc("/product/{id}", productcontroller.Show).Methods("GET")

	// detail product
	api.HandleFunc("/detail-product/{product_id}", detailproductcontroller.Show).Methods("GET")
	// show image
	api.HandleFunc("/product-image/{filename}", imagecontroller.Show).Methods("GET")

	// my cart
	api.HandleFunc("/my-cart", datausercontroller.ShowCart).Methods("GET")
	//cart
	api.HandleFunc("/cart", cartcontroller.Store).Methods("POST")
	api.HandleFunc("/cart/{id}", cartcontroller.Show).Methods("GET")
	api.HandleFunc("/cart/{id}", cartcontroller.Update).Methods("PUT")
	api.HandleFunc("/cart/{id}", cartcontroller.Delete).Methods("DELETE")

	api.Use(middlewares.JWTMiddleware)

	apiAdmin := r.PathPrefix("/api").Subrouter()
	// product
	apiAdmin.HandleFunc("/product", productcontroller.Store).Methods("POST")
	apiAdmin.HandleFunc("/product/{id}", productcontroller.Update).Methods("PUT")
	apiAdmin.HandleFunc("/product/{id}", productcontroller.Delete).Methods("DELETE")

	// detail product
	apiAdmin.HandleFunc("/detail-product", detailproductcontroller.Store).Methods("POST")
	apiAdmin.HandleFunc("/detail-product/{id}", detailproductcontroller.Update).Methods("PUT")
	apiAdmin.HandleFunc("/detail-product/{id}", detailproductcontroller.Delete).Methods("DELETE")

	apiAdmin.Use(middlewares.AdminRole)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4173"}, // Atur sesuai dengan asal lintas yang diizinkan
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Set-Cookie"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
