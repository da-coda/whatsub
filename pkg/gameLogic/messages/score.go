package messages

import "github.com/google/uuid"

const ScoreType Type = "score"

type Score struct {
	Envelop
	Payload struct {
		Scores []SingleScore
	}
}

type SingleScore struct {
	Name  string
	Score int
	UUID  uuid.UUID
}

func NewScoreMessage() Score {
	score := Score{}
	score.Envelop = Envelop{Type: ScoreType}
	score.Payload.Scores = make([]SingleScore, 0)
	return score

}
