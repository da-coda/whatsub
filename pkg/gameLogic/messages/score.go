package messages

const ScoreType = "score"

type Score struct {
	Envelop
	Payload struct {
		Score int
	}
}

func NewScoreMessage() Score {
	return Score{
		Envelop: Envelop{Type: ScoreType},
		Payload: struct {
			Score int
		}{},
	}
}
