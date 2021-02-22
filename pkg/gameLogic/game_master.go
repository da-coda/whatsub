package gameLogic

import (
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
	Worker         sync.Map
	workerShortIds sync.Map
}

//NewGameMaster creates a new GameMaster and starts the cleanUp routine for this game master
func NewGameMaster() *GameMaster {
	gm := &GameMaster{}
	go gm.cleanUp()
	return gm
}

//CreateGame is the main entrypoint for creating new games. It spawns a new Worker and puts the worker into the GameMaster.Worker map.
//Also starts the lobby for the created Worker and returns the UUID for the worker
func (gm *GameMaster) CreateGame() (uuid.UUID, string) {
	gameWorker := NewWorker()
	gm.Worker.Store(gameWorker.Id, gameWorker)
	gm.workerShortIds.Store(gameWorker.ShortId, gameWorker.Id)
	//OpenLobby on newly created game so that players can directly join
	go gameWorker.OpenLobby()
	return gameWorker.Id, gameWorker.ShortId
}

//JoinGame handles new clients joining a game by checking for worker by uuid, upgrading to websocket and creating a new Client
func (gm *GameMaster) JoinGame(w http.ResponseWriter, r *http.Request) {
	//get needed params from join request
	vars := mux.Vars(r)
	workerId := vars["uuid_or_key"]
	playerName := r.FormValue("name")

	//Parse uuid and find matching worker
	workerUuid, err := uuid.Parse(workerId)
	if err != nil {
		id, exists := gm.workerShortIds.Load(workerId)
		if !exists {
			logrus.WithError(err).Error("Unable to parse worker id")
			w.WriteHeader(400)
			return
		}
		workerUuid = id.(uuid.UUID)
	}

	gameWorker, exists := gm.Worker.Load(workerUuid)
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

	NewClient(c, playerName, gameWorker.(*Worker))
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
	gameWorker, exists := gm.Worker.Load(workerUuid)
	if !exists {
		logrus.WithField("UUID", workerUuid).Debug("Tried to join game on worker that does not exist")
		w.WriteHeader(404)
		return
	}
	go gameWorker.(*Worker).RunGame()
}

//cleanUp checks for every running worker if the worker is still needed.
//If not call Worker.Close and remove worker from GameMaster.Worker
func (gm *GameMaster) cleanUp() {
	for {
		gm.Worker.Range(func(key, value interface{}) bool {
			workerId := key.(uuid.UUID)
			gameWorker := value.(*Worker)
			if !gameWorker.StillNeeded() {
				logrus.WithField("Worker", workerId.String()).Debug("Removing worker in clean up")
				_ = gameWorker.Close()
				gm.Worker.Delete(workerId)
			}
			return true
		})
		time.Sleep(1 * time.Second)
	}
}
