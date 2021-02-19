package messages

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	HasNoType = errors.New("No type given in message")
	WrongType = errors.New("Wrong type")
)

type Type string

type Envelop struct {
	Type Type
}

func GetMessageType(jsonString []byte) (Type, error) {
	var answerMap map[string]json.RawMessage
	err := json.Unmarshal(jsonString, &answerMap)
	if err != nil {
		return "", err
	}
	var messageType string
	err = json.Unmarshal(answerMap["Type"], &messageType)
	if err != nil {
		return "", err
	}
	return Type(string(messageType)), nil
}
