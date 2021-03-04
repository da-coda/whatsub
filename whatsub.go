package main

import (
	"github.com/da-coda/whatsub/pkg/gameLogic"
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
	/*err := database.InitDB()
	if err != nil {
		logrus.WithError(err).Error("Unable to connect to DB")
		os.Exit(1)
	}*/
	gm := gameLogic.NewGameMaster()
	router := mux.NewRouter()
	router.HandleFunc("/game/create", gm.CreateGameHandler)
	//join
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("name", "{name}", "uuid", "{uuid}")
	//rejoin
	router.HandleFunc("/game/{uuid_or_key}/join", gm.JoinGame).Queries("uuid", "{uuid}")
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
