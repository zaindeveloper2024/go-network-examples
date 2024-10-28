package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	JWTSecret     string
	DatabaseURL   string
	Port          string
	TokenDuration time.Duration
}

var (
	config Config
)

func init() {
	config = Config{
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/auth_db?sslmode=disable"),
		Port:          getEnv("PORT", "8080"),
		TokenDuration: time.Hour * 24,
	}
}

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Status: status,
		Data:   data,
	}
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
	}
	jsonResponse(w, http.StatusOK, data)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health", healthHandler).Methods("GET")

	log.Printf("Server started on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
