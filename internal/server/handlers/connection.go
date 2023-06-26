package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type connectionHandler struct {
	logger *zap.Logger
}

func NewConnectionHandler(logger *zap.Logger) http.Handler {
	return &connectionHandler{
		logger: logger,
	}
}

func (h *connectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("failed to establish client connection", zap.Error(err))
	}
	h.logger.Info("new client connected", zap.Any("address", r.RemoteAddr))
	conn.Close()
}
