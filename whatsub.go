package main

import (
	"encoding/json"
	"github.com/da-coda/whatsub/pkg/gameLogic"
	"github.com/da-coda/whatsub/pkg/gameLogic/messages"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	gm := gameLogic.NewGameMaster()
	router := mux.NewRouter()
	router.HandleFunc("/game/create", CreateGameHandler(gm))
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("name", "{name}")
	router.HandleFunc("/game/{uuid}/start", gm.StartGame)
	router.PathPrefix("/").HandlerFunc(ServeWebpage).Methods("GET")
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
		uuid, key := gm.CreateGame()
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

func ServeWebpage(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join("src/", path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, "src/index.html")
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir("src")).ServeHTTP(w, r)
}
