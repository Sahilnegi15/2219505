package main

import (
	"log"
	"net/http"
)

func main() {
	router := SetupRoutes()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
