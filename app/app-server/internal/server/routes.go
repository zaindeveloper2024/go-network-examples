package server

import "net/http"

func (s *Server) setupRoutes() {
	// health
	s.router.HandleFunc("/health", s.handleHealth()).Methods(http.MethodGet)

	// users
	s.router.HandleFunc("/users", s.handleGetUsers()).Methods(http.MethodGet)
	s.router.HandleFunc("/users", s.handleCreateUser()).Methods(http.MethodPost)
	s.router.HandleFunc("/users/{id}", s.handleGetUser()).Methods(http.MethodGet)
	s.router.HandleFunc("/users/{id}", s.handleUpdateUser()).Methods(http.MethodPut)
	s.router.HandleFunc("/users/{id}", s.handleDeleteUser()).Methods(http.MethodDelete)
}
