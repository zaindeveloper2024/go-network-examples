package main

import (
	"fmt"
	"log"
	"net/http"

	"app-server/internal/config"
	"app-server/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	srv := server.NewServer(cfg)

	port := fmt.Sprintf(":%d", cfg.App.Port)

	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(port, srv.Router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
