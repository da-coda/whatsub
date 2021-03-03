package gameLogic

import (
	"bytes"
	"github.com/google/uuid"
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
	conn       *websocket.Conn
	uuid       uuid.UUID
	Name       string
	Worker     *Worker
	Send       chan []byte
	Message    chan []byte
	Blocked    bool
	close      chan bool
	log        *logrus.Entry
	Terminated bool
}

func NewClient(conn *websocket.Conn, name string, uuid uuid.UUID, gameWorker *Worker) *Client {
	client := &Client{
		conn:    conn,
		Name:    name,
		uuid:    uuid,
		Worker:  gameWorker,
		Send:    make(chan []byte, 256),
		Message: make(chan []byte),
		close:   make(chan bool),
	}
	client.log = logrus.WithField("Client", client.Name).WithField("Worker", gameWorker.Id.String())
	go client.readPump()
	go client.writePump()
	return client
}

func (c *Client) Close() error {
	c.log.Debug("Terminating client because Close() got called")
	c.Terminated = true
	c.close <- true
	close(c.Send)
	close(c.Message)
	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.log.Info("Closing websocket for read")
		c.Worker.Disconnect(c)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		c.log.WithError(err).Error("Unable to set read deadline")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		c.log.Trace("Pong")
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.WithError(err).Error("Client closed unexpected")
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if c.Blocked {
			c.log.WithField("Message", message).Trace("Message on blocked client")
			continue
		}
		if c.Worker.State != Closed {
			c.Message <- message
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.log.Info("Closing websocket for write")

		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case <-c.close:
			_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
			return
		case message, ok := <-c.Send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.log.WithError(err).Error("Unable to set write deadline")
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
				c.log.WithError(err).Error("Unable to write message")
			}
			if err := w.Close(); err != nil {
				if err != nil {
					c.log.WithError(err).Error("Unable to close writer")
				}
				return
			}
		case <-ticker.C:
			c.log.Trace("Ping")
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.log.WithError(err).Error("Unable to set write deadline")
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if err != nil {
					c.log.WithError(err).Info("Ping didn't reach client")
				}
				return
			}
		}
	}
}
