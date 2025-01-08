package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	userStore *UserStore
}

type UserStore struct {
	sync.RWMutex
	users map[string]User
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]User),
	}
}

func NewServer() *Server {
	s := &Server{
		router:    mux.NewRouter(),
		userStore: NewUserStore(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/health", s.handleHealth()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsers()).Methods("GET")
	s.router.HandleFunc("/users", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/users/{id}", s.handleGetUser()).Methods("GET")
}

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func (s *Server) handleUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.userStore.RLock()
		users := make([]User, 0, len(s.userStore.users))
		for _, user := range s.userStore.users {
			users = append(users, user)
		}
		s.userStore.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func (s *Server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user.ID = uuid.New().String()

		s.userStore.Lock()
		s.userStore.users[user.ID] = user
		s.userStore.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		s.userStore.RLock()
		user, exists := s.userStore.users[userID]
		s.userStore.RUnlock()

		if !exists {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func main() {
	server := NewServer()

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", server.router); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
