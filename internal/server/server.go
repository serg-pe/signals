package server

import (
	"github.com/serg-pe/signals/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
}

func New(cfg config.ServerConfig, logger *zap.Logger) (Server, error) {

	return Server{
		logger: logger,
	}, nil
}
