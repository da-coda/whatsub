package gameLogic

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

type GameMaster struct {
	Worker map[uuid.UUID]*Worker
}

func New() *GameMaster {
	gm := &GameMaster{}
	gm.Worker = make(map[uuid.UUID]*Worker)
	go gm.cleanUp()
	return gm
}

func (gm *GameMaster) CreateGame() uuid.UUID {
	gameWorker := NewWorker()
	gm.Worker[gameWorker.WorkerId] = gameWorker
	go gameWorker.OpenLobby()
	return gameWorker.WorkerId
}

func (gm *GameMaster) JoinGame(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Player joining a game")
	vars := mux.Vars(r)
	workerId := vars["uuid"]
	playerName := vars["name"]
	workerUuid, err := uuid.Parse(workerId)
	if err != nil {
		logrus.WithError(err).Error("Unable to parse worker id")
		w.WriteHeader(400)
		return
	}
	gameWorker, exists := gm.Worker[workerUuid]
	if !exists {
		logrus.WithField("UUID", workerUuid).Debug("Tried to join game on worker that does not exist")
		w.WriteHeader(404)
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Unable to upgrade to websocket connection")
		return
	}
	NewClient(c, playerName, gameWorker)
}

func (gm *GameMaster) StartGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workerId := vars["uuid"]
	workerUuid, err := uuid.Parse(workerId)
	if err != nil {
		logrus.WithError(err).Error("Unable to parse worker id")
		w.WriteHeader(400)
		return
	}
	gameWorker, exists := gm.Worker[workerUuid]
	if !exists {
		logrus.WithField("UUID", workerUuid).Debug("Tried to join game on worker that does not exist")
		w.WriteHeader(404)
		return
	}
	go gameWorker.RunGame()
}

func (gm *GameMaster) cleanUp() {
	for {
		for workerId, gameWorker := range gm.Worker {
			if !gameWorker.StillNeeded() {
				logrus.WithField("Worker", workerId.String()).Debug("Removing worker in clean up")
				delete(gm.Worker, workerId)
			}
		}
		time.Sleep(30 * time.Second)
	}
}
