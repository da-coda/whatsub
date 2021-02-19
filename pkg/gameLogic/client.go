package gameLogic

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn    *websocket.Conn
	Name    string
	Score   int
	Worker  *Worker
	Send    chan []byte
	Message chan []byte
	Blocked bool
}

func (c *Client) Close() error {
	logrus.Debug("Terminating client because Close() got called")
	close(c.Send)
	close(c.Message)
	_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	_ = c.conn.Close()
	return nil
}

func NewClient(conn *websocket.Conn, name string, gameWorker *Worker) {
	client := &Client{conn: conn, Name: name, Worker: gameWorker, Send: make(chan []byte, 256), Message: make(chan []byte)}
	gameWorker.Register <- client
	go client.readPump()
	go client.writePump()
}

func (c *Client) readPump() {
	defer func() {
		logrus.
			WithField("Worker", c.Worker.Id).
			WithField("Client", c.Name).
			Info("Closing websocket for read")
		if c.Worker.State != Closed {
			c.Worker.Unregister <- c
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logrus.WithError(err).
			WithField("Worker", c.Worker.Id).
			WithField("Client", c.Name).
			Error("Unable to set read deadline")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		logrus.WithField("Client", c.Name).WithField("Worker", c.Worker.Id).Trace("Pong")
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).
					WithField("Worker", c.Worker.Id).
					WithField("Client", c.Name).
					Error("Client closed unexpected")
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if c.Worker.State != Closed && !c.Blocked {
			c.Message <- message
			c.Worker.Incoming <- c
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logrus.
			WithField("Worker", c.Worker.Id).
			WithField("Client", c.Name).
			Info("Closing websocket for write")

		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logrus.WithError(err).
					WithField("Worker", c.Worker.Id).
					WithField("Client", c.Name).
					Error("Unable to set write deadline")
				return
			}
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				logrus.WithError(err).
					WithField("Worker", c.Worker.Id).
					WithField("Client", c.Name).
					Error("Unable to write message")
			}
			if err := w.Close(); err != nil {
				if err != nil {
					logrus.WithError(err).
						WithField("Worker", c.Worker.Id).
						WithField("Client", c.Name).
						Error("Unable to close writer")
				}
				return
			}
		case <-ticker.C:
			logrus.WithField("Client", c.Name).WithField("Worker", c.Worker.Id).Trace("Ping")
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logrus.WithError(err).
					WithField("Worker", c.Worker.Id).
					WithField("Client", c.Name).
					Error("Unable to set write deadline")
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if err != nil {
					logrus.WithError(err).
						WithField("Worker", c.Worker.Id).
						WithField("Client", c.Name).
						Info("Ping didn't reach client")
				}
				return
			}
		}
	}
}
