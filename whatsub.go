package main

import (
	"encoding/json"
	"github.com/da-coda/whatsub/pkg/gameMaster"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	gm := gameMaster.New()
	router := mux.NewRouter()
	router.Path("/startGame").HandlerFunc(CreateGameHandler(gm))
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func CreateGameHandler(gm *gameMaster.GameMaster) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		uuid := gm.StartGame()
		response := map[string]string{"uuid": uuid.String()}
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
