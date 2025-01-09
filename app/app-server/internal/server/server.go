package server

import (
	"app-server/internal/store"

	"github.com/gorilla/mux"
)

type Server struct {
	Router    *mux.Router
	userStore *store.UserStore
}

func NewServer() *Server {
	s := &Server{
		Router:    mux.NewRouter(),
		userStore: store.NewUserStore(),
	}
	s.setupRoutes()
	return s
}
