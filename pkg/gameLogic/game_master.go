package gameLogic

import (
	"crypto/md5"
	"encoding/json"
	"github.com/da-coda/whatsub/pkg/gameLogic/game"
	"github.com/da-coda/whatsub/pkg/gameLogic/messages"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"hash"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	MaxAllowedGamesPerIP = 5
	CleanUpInterval      = 1 * time.Second
)

var (
	TooManyOpenGames = errors.New("Client has to many open games")
)

//GameMaster is the central entrypoint and handler for all running games and creating games
type GameMaster struct {
	Worker                sync.Map
	workerShortIds        sync.Map
	hashedIpsRunningGames sync.Map
	factory               game.WorkerProvider
}

//NewGameMaster creates a new GameMaster and starts the cleanUp routine for this game master
func NewGameMaster(factory game.WorkerProvider) *GameMaster {
	gm := &GameMaster{factory: factory}
	go gm.cleanUp()
	return gm
}

//CreateGame is the main entrypoint for creating new games. It spawns a new Worker and puts the worker into the GameMaster.Worker map.
//Also starts the lobby for the created Worker and returns the UUID for the worker
func (gm *GameMaster) CreateGame(hashedIp hash.Hash, gameType string) (uuid.UUID, string, error) {
	sum := string(hashedIp.Sum(nil))
	gamesRunning, hasGamesRunning := gm.hashedIpsRunningGames.Load(sum)
	if hasGamesRunning && gamesRunning.(int) >= MaxAllowedGamesPerIP {
		return [16]byte{}, "", errors.Wrapf(TooManyOpenGames, "Client has already %d of %d games running.", gamesRunning.(int), MaxAllowedGamesPerIP)
	}
	if gamesRunning == nil {
		gamesRunning = 0
	}
	gameWorkerConst, err := gm.factory.GetConstructor(gameType)
	if err != nil {
		if errors.Is(err, game.UnknownGameTypeErr) {
			return [16]byte{}, "", errors.Wrapf(err, "Unknown game type given: %s", gameType)
		}
	}
	gameWorker := gameWorkerConst(sum)
	err = gameWorker.TransitionState(game.Open)
	if err != nil {
		return [16]byte{}, "", errors.Wrap(err, "Unable to transition worker to Open state")
	}

	gm.Worker.Store(gameWorker.ID(), gameWorker)
	gm.workerShortIds.Store(gameWorker.Key(), gameWorker.ID())
	gm.hashedIpsRunningGames.Store(sum, 1+gamesRunning.(int))
	return gameWorker.ID(), gameWorker.Key(), nil
}

//JoinGame handles new clients joining a game by checking for worker by uuid, upgrading to websocket and creating a new Client
func (gm *GameMaster) JoinGame(w http.ResponseWriter, r *http.Request) {
	//get needed params from join request
	vars := mux.Vars(r)
	workerId := vars["uuid_or_key"]
	gameWorker, err := gm.loadWorker(w, workerId)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		return
	}
	gameWorker.(game.Worker).Join(w, r)
}

//StartGame runs Worker.Run for given uuid and if worker exists
func (gm *GameMaster) StartGame(w http.ResponseWriter, r *http.Request) {
	//Get needed params from request
	vars := mux.Vars(r)
	workerId := vars["uuid_or_key"]

	gameWorker, err := gm.loadWorker(w, workerId)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		return
	}

	go gameWorker.(game.Worker).Start()

	w.WriteHeader(204)
}

func (gm *GameMaster) loadWorker(w http.ResponseWriter, workerId string) (game.Worker, error) {
	//Parse uuid and find matching worker
	workerUuid, err := uuid.Parse(workerId)
	if err != nil {
		id, exists := gm.workerShortIds.Load(workerId)
		if !exists {
			w.WriteHeader(400)
			return nil, errors.New("Unable to parse worker id")
		}
		workerUuid = id.(uuid.UUID)
	}

	//Find worker and start game
	gameWorker, exists := gm.Worker.Load(workerUuid)
	if !exists {
		errorMessage := "Tried to load game on worker that does not exist"
		logrus.WithField("UUID", workerUuid).Debug(errorMessage)
		w.WriteHeader(404)
		return nil, errors.New(errorMessage)
	}
	return gameWorker.(game.Worker), nil
}

func (gm *GameMaster) CreateGameHandler(writer http.ResponseWriter, request *http.Request) {
	hashedIp := md5.New()
	_, err := io.WriteString(hashedIp, request.RemoteAddr)
	if err != nil {
		logrus.WithError(err).Error("Unable to hash ip address")
		writer.WriteHeader(500)
		return
	}
	workerUuid, key, err := gm.CreateGame(hashedIp, "TopOfTheTop")
	if err != nil {
		logrus.WithError(err).Error("Unable to create worker")
		errMsg := messages.NewErrorMessage()
		errMsg.Payload.Error = errors.Cause(err).Error()
		payload, err := json.Marshal(errMsg)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		writer.WriteHeader(400)
		_, err = writer.Write(payload)
		if err != nil {
			logrus.WithError(err).Error("Unable to send error message")
			return
		}
		return
	}
	response := messages.NewCreatedGameMessage()
	response.Payload.UUID = workerUuid.String()
	response.Payload.Key = key
	payload, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(500)
		return
	}
	_, err = writer.Write(payload)
	if err != nil {
		writer.WriteHeader(500)
		return
	}
}

//GetStatus fetches the current status of an active game
func (gm *GameMaster) GetStatus(w http.ResponseWriter, r *http.Request) {
	//get needed params from join request
	vars := mux.Vars(r)
	workerId := vars["uuid_or_key"]
	gameWorker, err := gm.loadWorker(w, workerId)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		return
	}
	response := messages.NewGameStatusMessage(gameWorker.State().String())
	response.Payload.Player = gameWorker.Players()
	payload, err := json.Marshal(response)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		w.WriteHeader(500)
		return
	}
	_, err = w.Write(payload)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		w.WriteHeader(500)
		return
	}
}

//cleanUp checks for every running worker if the worker is still needed.
//If not call Worker.Close and remove worker from GameMaster.Worker
func (gm *GameMaster) cleanUp() {
	for {
		gm.Worker.Range(func(key, value interface{}) bool {
			workerId := key.(uuid.UUID)
			gameWorker := value.(game.Worker)
			if !gameWorker.StillNeeded() {
				hashedIp := gameWorker.Creator()
				gamesRunning, _ := gm.hashedIpsRunningGames.Load(hashedIp)
				gm.hashedIpsRunningGames.Store(hashedIp, 1+gamesRunning.(int))
				logrus.WithField("Worker", workerId.String()).Debug("Removing worker in clean up")
				_ = gameWorker.Close()
				gm.Worker.Delete(workerId)
			}
			return true
		})
		time.Sleep(CleanUpInterval)
	}
}
