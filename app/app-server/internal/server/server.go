package server

import (
	"app-server/internal/config"
	"app-server/internal/store"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	router    *mux.Router
	userStore *store.UserStore
	config    *config.Config
	db        *sqlx.DB
	// logger
}

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	s := &Server{
		router:    mux.NewRouter(),
		userStore: store.NewUserStore(),
		config:    cfg,
		db:        db,
		// logger
	}
	// setupMiddlewares
	s.setupRoutes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}
