package gameLogic

import (
	"github.com/da-coda/whatsub/pkg/gameLogic/game"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type WorkerMock struct {
	mock.Mock
}

func (w *WorkerMock) Close() error {
	args := w.Called()
	return args.Error(0)
}

func (w *WorkerMock) Join(_ http.ResponseWriter, _ *http.Request) {
	w.Called()
}

func (w *WorkerMock) Disconnect(gameClient *game.Client) {
	panic("implement me")
}

func (w *WorkerMock) Run() {
	panic("implement me")
}

func (w *WorkerMock) StillNeeded() bool {
	args := w.Called()
	return args.Bool(0)
}

func (w *WorkerMock) ID() uuid.UUID {
	args := w.Called()
	return args.Get(0).(uuid.UUID)
}

func (w *WorkerMock) State() game.State {
	args := w.Called()
	return args.Get(0).(game.State)
}

func (w *WorkerMock) Key() string {
	args := w.Called()
	return args.String(0)
}

func (w *WorkerMock) TransitionState(state game.State) error {
	args := w.Called(state)
	return args.Error(0)
}

func (w *WorkerMock) Creator() string {
	args := w.Called()
	return args.String(0)
}

type FactoryMock struct {
	mock.Mock
}

func (f *FactoryMock) GetConstructor(workerType string) (game.WorkerConstructor, error) {
	args := f.Called(workerType)
	return args.Get(0).(game.WorkerConstructor), args.Error(1)
}

func MockWorkerConstructor(returnMock game.Worker) game.WorkerConstructor {
	return func(ipHash string) game.Worker {
		return returnMock
	}
}
