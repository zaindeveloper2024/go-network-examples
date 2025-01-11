package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	httpServer := &http.Server{
		Addr:         port,
		Handler:      srv.Router,
		ReadTimeout:  time.Duration(cfg.App.ReadTimeout),
		WriteTimeout: time.Duration(cfg.App.WriteTimeout),
	}

	go func() {
		log.Printf("Server is running on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
