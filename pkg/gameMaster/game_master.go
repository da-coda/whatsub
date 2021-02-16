package gameMaster

import (
	"github.com/da-coda/whatsub/pkg/worker"
	"github.com/google/uuid"
)

type GameMaster struct {
	Worker map[uuid.UUID]*worker.Worker
}

func New() *GameMaster {
	gm := &GameMaster{}
	gm.Worker = make(map[uuid.UUID]*worker.Worker)
	return gm
}

func (gm GameMaster) StartGame() uuid.UUID {
	gameWorker := worker.New()
	gm.Worker[gameWorker.WorkerId] = gameWorker
	return gameWorker.WorkerId
}
