package main

import (
	"crud-go/config"
	"crud-go/controllers/authcontroller"
	"crud-go/controllers/categorycontroller"
	"crud-go/controllers/homecontroller"
	"crud-go/controllers/productcontroller"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	// Auth
	http.HandleFunc("/login", authcontroller.Index)
	http.HandleFunc("/register", authcontroller.IndexRegis)
	http.HandleFunc("/register/add", authcontroller.RegisterUser)

	// 1. Home
	http.HandleFunc("/", homecontroller.Welcome)

	// 2. Products
	http.HandleFunc("/products", productcontroller.Index)
	http.HandleFunc("/products/add", productcontroller.Add)
	http.HandleFunc("/products/detail", productcontroller.Detail)
	http.HandleFunc("/products/edit", productcontroller.Edit)
	http.HandleFunc("/products/delete", productcontroller.Delete)

	// 3. Categories
	http.HandleFunc("/categories", categorycontroller.Index)
	http.HandleFunc("/categories/add", categorycontroller.Add)
	http.HandleFunc("/categories/edit", categorycontroller.Edit)
	http.HandleFunc("/categories/delete", categorycontroller.Delete)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
