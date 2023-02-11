package conn

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"encoding/json"
	"errors"
)

var (
	ErrUndefinedAction = errors.New("undefined action type")
)

type TypeAction struct {
	Type string `json:"action"`
}

func (pc *PlayerConn) UnmarshalAction(msg []byte) (actions.IAction, error) {
	ta := &TypeAction{}
	if err := json.Unmarshal(msg, ta); err != nil {
		return nil, err
	}
	var a actions.IAction

	switch ta.Type {
	case string(actions.PingType):
		a = &actions.Ping{}
	case string(actions.ReadyType):
		a = &actions.Ready{}
	case string(actions.AnswerType):
		a = &actions.Answer{}
	case string(actions.CloseType):
		a = &actions.Close{}
	case string(actions.StartType):
		a = &actions.Start{}
	default:
		return nil, ErrUndefinedAction
	}

	if err := json.Unmarshal(msg, a); err != nil {
		return nil, err
	}

	return a, nil
}
