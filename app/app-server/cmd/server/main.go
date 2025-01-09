package main

import (
	"log"
	"net/http"

	"app-server/internal/server"
)

func main() {
	server := server.NewServer()

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", server.Router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
