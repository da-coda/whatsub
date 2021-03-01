package main

import (
	"crypto/md5"
	"encoding/json"
	"github.com/da-coda/whatsub/pkg/database"
	"github.com/da-coda/whatsub/pkg/gameLogic"
	"github.com/da-coda/whatsub/pkg/gameLogic/messages"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	err := database.InitDB()
	if err != nil {
		logrus.WithError(err).Error("Unable to connect to DB")
		os.Exit(1)
	}
	gm := gameLogic.NewGameMaster()
	router := mux.NewRouter()
	router.HandleFunc("/game/create", CreateGameHandler(gm))
	//join
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("name", "{name}", "uuid", "{uuid}")
	//rejoin
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("uuid", "{uuid}")
	router.HandleFunc("/game/{uuid}/start", gm.StartGame)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func CreateGameHandler(gm *gameLogic.GameMaster) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		hash := md5.New()
		_, err := io.WriteString(hash, request.RemoteAddr)
		if err != nil {
			logrus.WithError(err).Error("Unable to hash ip address")
		}
		uuid, key := gm.CreateGame(hash)
		response := messages.NewCreatedGameMessage()
		response.Payload.UUID = uuid.String()
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
}