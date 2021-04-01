package main

import (
	"github.com/da-coda/whatsub/pkg/gameLogic"
	"github.com/da-coda/whatsub/pkg/gameLogic/game"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	gm := gameLogic.NewGameMaster(game.WorkerFactory{})
	router := mux.NewRouter()
	registerRoutes(router, gm)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// registerRoutes sets up the REST endpoints for accessing the game
func registerRoutes(router *mux.Router, gm *gameLogic.GameMaster) {
	router.HandleFunc("/game/create", gm.CreateGameHandler).Methods("POST")
	//join as new User
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("name", "{name}", "uuid", "{uuid}")
	//re-join or lobby
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("uuid", "{uuid}")
	//request status
	router.HandleFunc("/game/{uuid_or_key}/status", gm.GetStatus).Queries("uuid", "{uuid}")
	//start the game
	router.HandleFunc("/game/{uuid_or_key}/start", gm.StartGame)
}
