package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func jsonResponse(w http.ResponseWriter, status int, reponse Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(reponse)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Server is healthy",
		Data: map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
	jsonResponse(w, http.StatusOK, response)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		response := Response{
			Status:  "success",
			Message: "Hello, world!",
			Data: map[string]string{
				"version": "1.0.0",
			},
		}
		jsonResponse(w, http.StatusOK, response)
	default:
		jsonResponse(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Method not allowed",
		})
	}
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(startTime))
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				jsonResponse(w, http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Internal server error",
					Data: map[string]string{
						"message": fmt.Sprintf("%v", err),
						"code":    "500",
					},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func getServerConfig() ServerConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return ServerConfig{
		Port:         ":" + "8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	config := getServerConfig()

	// mux router
	router := mux.NewRouter()

	router.Use(logMiddleware)
	router.Use(recoverMiddleware)

	router.HandleFunc("/", helloHandler)
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	server := &http.Server{
		Handler:      router,
		Addr:         config.Port,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	log.Printf("Server listening on port %s", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
