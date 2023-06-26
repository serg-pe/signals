package server

import (
	"fmt"
	"net/http"

	"github.com/serg-pe/signals/internal/config"
	"github.com/serg-pe/signals/internal/server/handlers"
	"go.uber.org/zap"
)

type Server struct {
	logger   *zap.Logger
	cfg      config.ServerConfig
	serveMux *http.ServeMux
}

func New(cfg config.ServerConfig, logger *zap.Logger) (Server, error) {
	server := Server{
		logger:   logger,
		cfg:      cfg,
		serveMux: http.NewServeMux(),
	}

	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	s.serveMux.Handle("/connection/", handlers.NewConnectionHandler(s.logger))
}

func (s *Server) Run() error {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.cfg.Ip, s.cfg.Port), s.serveMux)
	if err != nil {
		return err
	}

	return nil
}
