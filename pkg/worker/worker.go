package worker

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Worker struct {
	WorkerId    uuid.UUID
	Connections []*websocket.Conn
}

func New() *Worker {
	return &Worker{WorkerId: uuid.New()}
}

func (worker *Worker) AddPlayer(conn *websocket.Conn) {
	worker.Connections = append(worker.Connections, conn)
}

func (worker Worker) RunGame() {
	fmt.Printf("Worker %s started a game", worker.WorkerId.String())
}
