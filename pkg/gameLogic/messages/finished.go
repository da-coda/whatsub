package messages

const FinishedType Type = "finished"

type Finished struct {
	Envelop
	Payload struct {
		Scores map[string]int
	}
}

func NewFinishedMessage() Finished {
	finished := Finished{}
	finished.Envelop = Envelop{Type: FinishedType}
	finished.Payload.Scores = make(map[string]int)
	return finished

}
