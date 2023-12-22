package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

	wg    *sync.WaitGroup
	pubMu *sync.Mutex
	pub   *websocket.Conn
	subMu *sync.Mutex
	sub   *websocket.Conn
}

func New(cfg config.ServerConfig, logger *zap.Logger) (Server, error) {
	server := Server{
		logger: logger.Named("server"),
		cfg:    cfg,

		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferPool: &sync.Pool{},
			CheckOrigin:     func(r *http.Request) bool { return true },
		},

		wg:    &sync.WaitGroup{},
		pubMu: &sync.Mutex{},
		subMu: &sync.Mutex{},
	}

	return server, nil
}

func (s *Server) setupRoutes() *http.ServeMux {
	// s.serveMux.Handle("/connection/", handlers.NewConnectionHandler(s.logger))
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
		s.sub = conn
		s.logger.Info("subscriber connected", zap.String("address", r.RemoteAddr))
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for {
				msgType, msg, err := s.sub.ReadMessage()
				if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					break
				}
				if err != nil {
					s.logger.Error("read message", zap.Error(err))
					break
				}
				s.logger.Debug("new message", zap.Int("type", msgType), zap.String("msg", hex.EncodeToString(msg)))
			}
			s.logger.Info("client disconnected")
			s.sub = nil
		}()
	}

	if s.pub != nil {

	}
}

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.cfg.Ip, s.cfg.Port),
		Handler:      s.setupRoutes(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	s.server = server

	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	if s.sub != nil {
		s.subMu.Lock()
		defer s.subMu.Unlock()
		err := s.sub.Close()
		if err != nil {
			s.logger.Debug("close connection error", zap.String("subscriber", s.sub.RemoteAddr().String()))
		}
	}
	if s.pub != nil {
		s.pubMu.Lock()
		defer s.pubMu.Unlock()
		err := s.pub.Close()
		if err != nil {
			s.logger.Debug("close connection error", zap.String("publisher", s.pub.RemoteAddr().String()))
		}
	}

	if s.server != nil {
		s.server.Shutdown(ctx)
	}

	s.wg.Wait()
}
