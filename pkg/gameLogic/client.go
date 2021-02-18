package gameLogic

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Client struct {
	*websocket.Conn
	*sync.RWMutex
	Name   string
	Score  int
	Worker *Worker
}

func NewClient(conn *websocket.Conn, name string, gameWorker *Worker) {
	client := &Client{Conn: conn, Name: name, Worker: gameWorker}
	gameWorker.Register <- client
}

func (c *Client) readPump() {
	defer func() {
		c.Worker.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}
