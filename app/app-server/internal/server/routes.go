package server

import "net/http"

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("/health", s.handleHealth()).Methods(http.MethodGet)
	s.Router.HandleFunc("/users", s.handleGetUsers()).Methods(http.MethodGet)
	s.Router.HandleFunc("/users", s.handleCreateUser()).Methods(http.MethodPost)
	s.Router.HandleFunc("/users/{id}", s.handleGetUser()).Methods(http.MethodGet)
	s.Router.HandleFunc("/users/{id}", s.handleUpdateUser()).Methods(http.MethodPut)
	s.Router.HandleFunc("/users/{id}", s.handleDeleteUser()).Methods(http.MethodDelete)
}
