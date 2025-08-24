package main

import (
	"crud-go/config"
	"crud-go/controllers/authcontroller"
	"crud-go/controllers/categorycontroller"
	"crud-go/controllers/homecontroller"
	"crud-go/controllers/productcontroller"
	"crud-go/middlewares"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	// Auth routes (guest only)
	http.HandleFunc("/login", middlewares.GuestOnly(authcontroller.Index))
	http.HandleFunc("/login/auth", middlewares.GuestOnly(authcontroller.LoginUser))
	http.HandleFunc("/register", middlewares.GuestOnly(authcontroller.IndexRegis))
	http.HandleFunc("/register/add", middlewares.GuestOnly(authcontroller.RegisterUser))
	
	// Logout route (authenticated users only)
	http.HandleFunc("/logout", middlewares.AuthRequired(authcontroller.LogoutUser))

	// 1. Home (protected)
	http.HandleFunc("/", middlewares.AuthRequired(homecontroller.Welcome))

	// 2. Products (protected)
	http.HandleFunc("/products", middlewares.AuthRequired(productcontroller.Index))
	http.HandleFunc("/products/add", middlewares.AuthRequired(productcontroller.Add))
	http.HandleFunc("/products/detail", middlewares.AuthRequired(productcontroller.Detail))
	http.HandleFunc("/products/edit", middlewares.AuthRequired(productcontroller.Edit))
	http.HandleFunc("/products/delete", middlewares.AuthRequired(productcontroller.Delete))

	// 3. Categories (protected)
	http.HandleFunc("/categories", middlewares.AuthRequired(categorycontroller.Index))
	http.HandleFunc("/categories/add", middlewares.AuthRequired(categorycontroller.Add))
	http.HandleFunc("/categories/edit", middlewares.AuthRequired(categorycontroller.Edit))
	http.HandleFunc("/categories/delete", middlewares.AuthRequired(categorycontroller.Delete))

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
