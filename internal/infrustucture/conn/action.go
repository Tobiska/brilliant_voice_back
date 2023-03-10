package conn

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"encoding/json"
	"errors"
)

var (
	ErrUndefinedAction     = errors.New("undefined action type")
	ErrDoesNotContainsType = errors.New("msg doesn't contains type")
	ErrDoesNotContainsText = errors.New("msg answer doesn't contains text")
)

type TypeAction struct {
	Type string `json:"action"`
}

func (pc *PlayerConn) UnmarshalAction(msg []byte) (fsm.IUserAction, error) {
	var bt map[string]string
	if err := json.Unmarshal(msg, &bt); err != nil {
		return nil, err
	}
	action, ex := bt["action"]
	if !ex {
		return nil, ErrDoesNotContainsType
	}
	var a fsm.IUserAction

	switch action {
	case "answer":
		text, ex := bt["text"]
		if !ex {
			return nil, ErrDoesNotContainsText
		}
		a = actions.AnswerAction(pc.user, text)
	case "start":
		a = actions.StartAction(pc.user)
	default:
		return nil, ErrUndefinedAction
	}

	return a, nil
}
