package main

import (
	"crud-go/config"
	"crud-go/controllers/homecontroller"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	// 1. Home
	http.HandleFunc("/", homecontroller.Welcome)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
