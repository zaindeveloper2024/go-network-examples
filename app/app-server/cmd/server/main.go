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

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"app-server/internal/config"
	"app-server/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	db, err := sqlx.Connect("postgres", cfg.DB.DNS())
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer db.Close()

	srv := server.NewServer(cfg, db)
	port := fmt.Sprintf(":%d", cfg.App.Port)

	httpServer := newHTTPServer(srv, port, cfg)

	return serveHTTP(httpServer, port, cfg)
}

func newHTTPServer(srv *server.Server, port string, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         port,
		Handler:      srv.Handler(),
		ReadTimeout:  time.Duration(cfg.App.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.App.WriteTimeout) * time.Second,
	}
}

func serveHTTP(httpServer *http.Server, port string, cfg *config.Config) error {
	errChan := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server is running on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("could not start server: %v", err)
	case sig := <-shutdown:
		return gracefulShutdown(httpServer, cfg, sig)
	}
}

func gracefulShutdown(httpServer *http.Server, cfg *config.Config, sig os.Signal) error {
	log.Printf("Received signal %v, initiating graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.App.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}

	log.Println("Server shutdown completed")
	return nil
}
