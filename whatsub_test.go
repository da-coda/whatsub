package main

import (
	"encoding/json"
	"github.com/da-coda/whatsub/pkg/gameLogic"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_createGame(t *testing.T) {
	gm := gameLogic.NewGameMaster()
	router := mux.NewRouter()
	registerRoutes(router, gm)

	type responseFormat struct {
		Type    string `json:"Type"`
		Payload struct {
			UUID string `json:"UUID"`
			Key  string `json:"Key"`
		} `json:"Payload"`
	}

	tests := []struct {
		name                 string
		method               string
		expectedResponseCode int
	}{
		{
			name:                 "Create a new Game",
			method:               "POST",
			expectedResponseCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, "/game/create", nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedResponseCode {
				t.Errorf("Got %d, Want %d", rr.Code, tt.expectedResponseCode)
			}
			if tt.expectedResponseCode != 200 {
				return
			}
			var unmarshalledResponse responseFormat
			err := json.Unmarshal(rr.Body.Bytes(), &unmarshalledResponse)
			if err != nil {
				t.Errorf("Could not unmarshall %s", rr.Body.String())
			}
			if unmarshalledResponse.Type != "created_game" {
				t.Errorf("Got %s, Want 'created_game'", unmarshalledResponse.Type)
			}
		})
	}
}
