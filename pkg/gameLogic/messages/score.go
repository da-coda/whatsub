package messages

const ScoreType Type = "score"

type Score struct {
	Envelop
	Payload struct {
		Scores map[string]int
	}
}

func NewScoreMessage() Score {
	score := Score{}
	score.Envelop = Envelop{Type: ScoreType}
	score.Payload.Scores = make(map[string]int)
	return score

}
