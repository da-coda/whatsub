package gameLogic

import (
	"crypto/md5"
	"encoding/json"
	"github.com/da-coda/whatsub/pkg/gameLogic/game"
	"github.com/da-coda/whatsub/pkg/gameLogic/messages"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type GameMasterTestSuite struct {
	suite.Suite
}

func (suite *GameMasterTestSuite) TestNewGameMaster_StartCleanup() {
	workerMock := new(WorkerMock)
	factoryMock := new(FactoryMock)
	workerMock.On("Creator").Return("test")
	workerMock.On("Close").Return(nil)
	workerMock.On("StillNeeded").Return(false)
	tmpUUID := uuid.New()
	gm := NewGameMaster(factoryMock)
	gm.Worker.Store(tmpUUID, workerMock)
	gm.hashedIpsRunningGames.Store("test", 1)
	_, hasWorker := gm.Worker.Load(tmpUUID)
	suite.True(hasWorker)
	//assure that clean up has run atleast once
	time.Sleep(CleanUpInterval + 500*time.Millisecond)
	_, hasWorker = gm.Worker.Load(tmpUUID)
	suite.False(hasWorker)
	workerMock.AssertExpectations(suite.T())
	factoryMock.AssertExpectations(suite.T())
}

func (suite *GameMasterTestSuite) TestCreateGame() {
	fakeUUID := uuid.New()
	workerMock := new(WorkerMock)
	workerMock.On("TransitionState", game.Open).Return(nil)
	workerMock.On("StillNeeded").Maybe().Return(true)
	workerMock.On("ID").Return(fakeUUID)
	workerMock.On("Key").Return("testKey")

	factoryMock := new(FactoryMock)
	factoryMock.On("GetConstructor", "test").Return(MockWorkerConstructor(workerMock), nil)

	fakeHash := md5.New()
	_, err := io.WriteString(fakeHash, "127.0.0.1")
	suite.NoError(err)
	gm := NewGameMaster(factoryMock)
	id, key, err := gm.CreateGame(fakeHash, "test")
	suite.NoError(err)
	suite.NotEmpty(key)
	_, hasWorker := gm.Worker.Load(id)
	suite.True(hasWorker)
	workerMock.AssertExpectations(suite.T())
	factoryMock.AssertExpectations(suite.T())
}

func (suite *GameMasterTestSuite) TestCreateGame_MaxGames() {

	workerMock := new(WorkerMock)
	workerMock.On("TransitionState", game.Open).Return(nil)
	workerMock.On("Key").Return("testKey")
	workerMock.On("StillNeeded").Maybe().Return(true)
	factoryMock := new(FactoryMock)
	factoryMock.On("GetConstructor", "test").Return(MockWorkerConstructor(workerMock), nil)
	fakeHash := md5.New()
	_, err := io.WriteString(fakeHash, "127.0.0.1")
	suite.NoError(err)
	gm := NewGameMaster(factoryMock)

	//These should all be created
	for i := 0; i < MaxAllowedGamesPerIP; i++ {
		workerMock.On("ID").Return(uuid.New())
		id, key, err := gm.CreateGame(fakeHash, "test")
		suite.NoError(err)
		suite.NotEmpty(key)
		_, hasWorker := gm.Worker.Load(id)
		suite.True(hasWorker)
	}

	//This one should be not created
	workerMock.On("ID").Return(uuid.New())
	_, _, err = gm.CreateGame(fakeHash, "test")
	suite.ErrorIs(err, TooManyOpenGames)
	workerMock.AssertExpectations(suite.T())
	factoryMock.AssertExpectations(suite.T())
}

func (suite *GameMasterTestSuite) TestCreateGameHandler() {
	request, _ := http.NewRequest(http.MethodPost, "/game/create", nil)
	request.RemoteAddr = "127.0.0.1"
	response := httptest.NewRecorder()

	u := uuid.New()
	key := "testKey"
	workerMock := new(WorkerMock)
	workerMock.On("TransitionState", game.Open).Return(nil)
	workerMock.On("Key").Return(key)
	workerMock.On("StillNeeded").Maybe().Return(true)
	workerMock.On("ID").Return(u)
	factoryMock := new(FactoryMock)
	factoryMock.On("GetConstructor", "TopOfTheTop").Return(MockWorkerConstructor(workerMock), nil)
	gm := NewGameMaster(factoryMock)
	gm.CreateGameHandler(response, request)
	suite.Equal(200, response.Code)
	body, err := io.ReadAll(response.Body)
	suite.NoError(err)
	var msg messages.CreatedGame
	err = json.Unmarshal(body, &msg)
	suite.NoError(err)
	suite.Equal(u.String(), msg.Payload.UUID)
	suite.Equal(key, msg.Payload.Key)
}

func (suite *GameMasterTestSuite) TestJoinGame() {

}

func TestGameMaster(t *testing.T) {
	suite.Run(t, new(GameMasterTestSuite))
}
