package messages

import "github.com/google/uuid"

const JoinedType Type = "joined"

type Joined struct {
	Envelop
	Payload struct {
		Name string
		UUID uuid.UUID
	}
}

func UserJoinedMessage() Joined {
	return Joined{
		Envelop: Envelop{Type: JoinedType},
		Payload: struct {
			Name string
			UUID uuid.UUID
		}{},
	}
}
