package gameLogic

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

//GameMaster is the central entrypoint and handler for all running games and creating games
type GameMaster struct {
	Worker map[uuid.UUID]*Worker
	mu     *sync.Mutex
}

//NewGameMaster creates a new GameMaster and starts the cleanUp routine for this game master
func NewGameMaster() *GameMaster {
	gm := &GameMaster{}
	gm.Worker = make(map[uuid.UUID]*Worker)
	gm.mu = new(sync.Mutex)
	//start GameMaster.cleanUp routine for removing unneeded worker
	go gm.cleanUp()
	return gm
}

//CreateGame is the main entrypoint for creating new games. It spawns a new Worker and puts the worker into the GameMaster.Worker map.
//Also starts the lobby for the created Worker and returns the UUID for the worker
func (gm *GameMaster) CreateGame() uuid.UUID {
	gameWorker := NewWorker()
	gm.mu.Lock()
	gm.Worker[gameWorker.Id] = gameWorker
	gm.mu.Unlock()
	//OpenLobby on newly created game so that players can directly join
	go gameWorker.OpenLobby()
	return gameWorker.Id
}

//JoinGame handles new clients joining a game by checking for worker by uuid, upgrading to websocket and creating a new Client
func (gm *GameMaster) JoinGame(w http.ResponseWriter, r *http.Request) {
	//get needed params from join request
	vars := mux.Vars(r)
	workerId := vars["uuid"]
	playerName := r.FormValue("name")

	//Parse uuid and find matching worker
	workerUuid, err := uuid.Parse(workerId)
	if err != nil {
		logrus.WithError(err).Error("Unable to parse worker id")
		w.WriteHeader(400)
		return
	}
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gameWorker, exists := gm.Worker[workerUuid]
	if !exists {
		logrus.WithField("UUID", workerUuid).Debug("Tried to join game on worker that does not exist")
		w.WriteHeader(404)
		return
	}

	//Upgrade connection to websocket and create new client
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Unable to upgrade to websocket connection")
		return
	}
	NewClient(c, playerName, gameWorker)
}

//StartGame runs Worker.RunGame for given uuid and if worker exists
func (gm *GameMaster) StartGame(w http.ResponseWriter, r *http.Request) {
	//Get needed params from request
	vars := mux.Vars(r)
	workerId := vars["uuid"]
	workerUuid, err := uuid.Parse(workerId)
	if err != nil {
		logrus.WithError(err).Error("Unable to parse worker id")
		w.WriteHeader(400)
		return
	}

	//Find worker and start game
	gm.mu.Lock()
	defer gm.mu.Unlock()
	fmt.Println(len(gm.Worker))
	gameWorker, exists := gm.Worker[workerUuid]
	if !exists {
		logrus.WithField("UUID", workerUuid).Debug("Tried to join game on worker that does not exist")
		w.WriteHeader(404)
		return
	}
	go gameWorker.RunGame()
}

//cleanUp checks for every running worker if the worker is still needed.
//If not call Worker.Close and remove worker from GameMaster.Worker
func (gm *GameMaster) cleanUp() {
	for {
		for workerId, gameWorker := range gm.Worker {
			if !gameWorker.StillNeeded() {
				logrus.WithField("Worker", workerId.String()).Debug("Removing worker in clean up")
				_ = gameWorker.Close()
				delete(gm.Worker, workerId)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
