package main

import (
	"crud-go/config"
	"crud-go/controllers/categorycontroller"
	"crud-go/controllers/homecontroller"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	// 1. Home
	http.HandleFunc("/", homecontroller.Welcome)

	// 2. Categories
	http.HandleFunc("/categories", categorycontroller.Index)
	http.HandleFunc("/categories/add", categorycontroller.Add)
	http.HandleFunc("/categories/edit", categorycontroller.Edit)
	http.HandleFunc("/categories/delete", categorycontroller.Delete)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
