package game

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

//go:generate stringer -type=State
type State int
type WorkerConstructor = func(ipHash string) Worker

var (
	UnknownGameTypeErr        = errors.New("game type is not supported")
	IllegalStateTransitionErr = errors.New("state transition is not possible")
)

const (
	Created State = iota
	Open
	Started
	Done
	Closed
)

var (
	AllowedTransitions = map[State][]State{
		Created: {Open, Closed},
		Open:    {Started, Done, Closed},
		Started: {Done, Closed},
		Done:    {Closed},
		Closed:  {Done},
	}
)

type Worker interface {
	io.Closer
	Join(w http.ResponseWriter, r *http.Request)
	Disconnect(gameClient *Client)
	Start()
	Run()
	StillNeeded() bool
	ID() uuid.UUID
	State() State
	Key() string
	TransitionState(state State) error
	Creator() string
	Players() []string
}

func CanTransition(currentState State, transitionState State) bool {
	possibleTransitions := AllowedTransitions[currentState]
	for _, transState := range possibleTransitions {
		if transState == transitionState {
			return true
		}
	}
	return false
}
