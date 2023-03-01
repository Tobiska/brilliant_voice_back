package conn

import (
	"brillian_voice_back/internal/domain/entity/actions"
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
	case "answer":
		a = actions.Answer{}
	case "start":
		a = actions.Start{}
	default:
		return nil, ErrUndefinedAction
	}

	if err := json.Unmarshal(msg, a); err != nil {
		return nil, err
	}

	return a, nil
}
