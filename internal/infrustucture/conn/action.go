package conn

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"encoding/json"
	"errors"
)

var (
	ErrUndefinedAction = errors.New("undefined action type")
)

type TypeAction struct {
	Type string `json:"action"`
}

func (pc *PlayerConn) UnmarshalAction(msg []byte) (fsm.IAction, error) {
	ta := &TypeAction{}
	if err := json.Unmarshal(msg, ta); err != nil {
		return nil, err
	}
	var a fsm.IAction

	switch ta.Type {
	case string(fsm.PingType):
		a = &fsm.Ping{}
	case string(fsm.ReadyType):
		a = &fsm.Ready{}
	case string(fsm.AnswerType):
		a = &fsm.Answer{}
	case string(fsm.CloseType):
		a = &fsm.Close{}
	case string(fsm.StartType):
		a = &fsm.Start{}
	default:
		return nil, ErrUndefinedAction
	}

	if err := json.Unmarshal(msg, a); err != nil {
		return nil, err
	}

	return a, nil
}
