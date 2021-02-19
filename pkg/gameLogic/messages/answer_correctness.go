package messages

const AnswerCorrectnessType = "answer_correct"

type AnswerCorrectness struct {
	Envelop
	Payload struct {
		Correct bool
	}
}

func NewAnswerCorrectnessMessage() AnswerCorrectness {
	return AnswerCorrectness{
		Envelop: Envelop{Type: AnswerCorrectnessType},
		Payload: struct {
			Correct bool
		}{},
	}
}
