package main

import (
	"github.com/da-coda/whatsub/pkg/gameLogic"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	/*err := database.InitDB()
	if err != nil {
		logrus.WithError(err).Error("Unable to connect to DB")
		os.Exit(1)
	}*/
	gm := gameLogic.NewGameMaster()
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
	router.HandleFunc("/game/create", gm.CreateGameHandler)
	//join as new User
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("name", "{name}", "uuid", "{uuid}")
	//re-join or lobby
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("uuid", "{uuid}")
	router.HandleFunc("/game/{uuid}/start", gm.StartGame)
}
