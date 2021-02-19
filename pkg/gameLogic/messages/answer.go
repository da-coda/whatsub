package messages

import "encoding/json"

const AnswerType Type = "answer"

type Answer struct {
	Envelop
	Payload struct {
		Answer string
	}
}

func (answer *Answer) Parse(jsonString []byte) error {
	messageType, err := GetMessageType(jsonString)
	if err != nil {
		return err
	}
	if messageType != AnswerType {
		return WrongType
	}
	err = json.Unmarshal(jsonString, answer)
	return err
}
