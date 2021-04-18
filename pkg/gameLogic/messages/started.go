package messages

const StartedType Type = "started"

type Started struct {
	Envelop
}

func NewStartedMessage() Started {
	return Started{
		Envelop: Envelop{Type: StartedType},
	}
}