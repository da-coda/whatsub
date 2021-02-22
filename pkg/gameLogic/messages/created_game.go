package messages

const CreatedGameType Type = "created_game"

type CreatedGame struct {
	Envelop
	Payload struct {
		UUID string
		Key  string
	}
}

func NewCreatedGameMessage() CreatedGame {
	return CreatedGame{
		Envelop: Envelop{Type: CreatedGameType},
		Payload: struct {
			UUID string
			Key  string
		}{},
	}
}
