package messages

const ErrorType Type = "error"

type Error struct {
	Envelop
	Payload struct {
		Error string
	}
}

func NewErrorMessage() Error {
	errorMsg := Error{}
	errorMsg.Envelop = Envelop{Type: ErrorType}
	return errorMsg
}
