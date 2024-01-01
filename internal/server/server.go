package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/serg-pe/signals/internal/client"
	"github.com/serg-pe/signals/internal/config"
	"go.uber.org/zap"
)

const (
	queryIsPublisherName = "is-initiator"
)

type Server struct {
	logger *zap.Logger
	cfg    config.ServerConfig

	server *http.Server

	upgrader websocket.Upgrader

	client *client.Client

	wg *sync.WaitGroup
}

func New(cfg config.ServerConfig, logger *zap.Logger) (Server, error) {
	s := Server{
		logger: logger.Named("server"),
		cfg:    cfg,

		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Ip, cfg.Port),
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		},

		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferPool: &sync.Pool{},
			CheckOrigin:     func(r *http.Request) bool { return true },
		},

		wg: &sync.WaitGroup{},
	}

	return s, nil
}

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/connection/", s.connect)

	return mux
}

func (s *Server) connect(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		isPub bool
	)

	if !websocket.IsWebSocketUpgrade(r) {
		s.logger.Debug("no upgrade protocol header", zap.String("from", r.RemoteAddr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.URL.Query().Has(queryIsPublisherName) {
		isPubRaw := r.URL.Query().Get(queryIsPublisherName)
		isPub, err = strconv.ParseBool(isPubRaw)
		if err != nil {
			s.logger.Debug("parse bool error", zap.String("client", r.RemoteAddr), zap.String("value", isPubRaw))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("upgrade connection", zap.String("client", r.RemoteAddr), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isPub {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.logger.Info("subscriber connected", zap.String("address", conn.RemoteAddr().String()))
			s.client = client.New(s.logger.Named(fmt.Sprintf("client %s", conn.RemoteAddr().String())), conn)
			go s.client.Listen()
		}()
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.server.Handler = s.setupRoutes()

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Debug("shutdown error", zap.Error(err))
	}

	if s.client != nil {
		s.client.Stop()
	}
	s.wg.Wait()

	s.logger.Info("server stopped")
}
