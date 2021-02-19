package messages

const AnswerCorrectnessType Type = "answer_correct"

type AnswerCorrectness struct {
	Envelop
	Payload struct {
		Correct       bool
		CorrectAnswer string
	}
}

func NewAnswerCorrectnessMessage() AnswerCorrectness {
	return AnswerCorrectness{
		Envelop: Envelop{Type: AnswerCorrectnessType},
		Payload: struct {
			Correct       bool
			CorrectAnswer string
		}{},
	}
}
