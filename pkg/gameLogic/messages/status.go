package messages

const GameStatusType Type = "game_status"

type GameStatus struct {
	Envelop
	Payload struct {
		State  string
		Player []string
	}
}

func NewGameStatusMessage(state string) GameStatus {
	return GameStatus{
		Envelop: Envelop{Type: GameStatusType},
		Payload: struct {
			State  string
			Player []string
		}{
			state,
			make([]string, 1),
		},
	}
}
