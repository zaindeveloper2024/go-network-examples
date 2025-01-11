package server

import (
	"app-server/internal/config"
	"app-server/internal/store"

	"github.com/gorilla/mux"
)

type Server struct {
	Router    *mux.Router
	userStore *store.UserStore
	config    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		Router:    mux.NewRouter(),
		userStore: store.NewUserStore(),
		config:    cfg,
	}
	s.setupRoutes()
	return s
}
