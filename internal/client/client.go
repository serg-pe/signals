package client

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/serg-pe/signals/pkg/signals"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	conn   *websocket.Conn
	stop   chan struct{}
}

func New(logger *zap.Logger, conn *websocket.Conn) *Client {
	return &Client{
		logger: logger,
		conn:   conn,
		stop:   make(chan struct{}),
	}
}

func (c *Client) Listen() {
	defer c.close()

	msgBuf := make([]byte, 1)

	for {
		select {
		case <-c.stop:
			return
		default:
			msgType, msg, err := c.conn.ReadMessage()
			if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.logger.Debug("client closed connection", zap.Error(err))
				return
			}
			if err != nil {
				c.logger.Error("listen client error", zap.Error(err))
				return
			}

			if msgType != websocket.BinaryMessage {
				c.logger.Debug("non binary msg type", zap.Int("type", msgType))
				continue
			}

			switch signals.Signal(msg[0]) {
			case signals.SignalPing:
				msgBuf[0] = byte(signals.SignalPong)
				err = c.conn.WriteMessage(websocket.BinaryMessage, msgBuf)
			default:
				c.logger.Debug("got message", zap.Int8("msg", int8(msg[0])))
			}
			if err != nil {
				c.logger.Debug("reply client error", zap.Error(err))
				return
			}
		}
	}
}

func (c *Client) Stop() {
	c.stop <- struct{}{}
}

func (c *Client) close() {
	time.Sleep(time.Second * 5)
	err := c.conn.Close()
	if err != nil {
		c.logger.Error("close error", zap.Error(err))
	}
	c.conn = nil
	c.logger.Info("client disconnected")
}
