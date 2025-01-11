package server

import (
	"app-server/internal/config"
	"app-server/internal/store"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	userStore *store.UserStore
	config    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		router:    mux.NewRouter(),
		userStore: store.NewUserStore(),
		config:    cfg,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}
