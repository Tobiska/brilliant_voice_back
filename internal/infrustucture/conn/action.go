package conn

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	ErrUndefinedAction  = errors.New("undefined action type")
	ErrTypeDoesNotExist = errors.New("action type doesn't exist")
)

type IActionInf interface {
	ToDomain(u *game.User) (fsm.IUserAction, error)
}

type TypeAction struct {
	Type string `json:"action"`
}

type Answer struct {
	TypeAction
	Text string `json:"text"`
}

func (a Answer) ToDomain(u *game.User) (fsm.IUserAction, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return actions.AnswerAction(u, a.Text), nil
}

func (a Answer) validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Text, validation.NotNil, validation.Length(3, 30)),
	)
}

type Start struct {
	TypeAction
}

func (a Start) ToDomain(u *game.User) (fsm.IUserAction, error) {
	return actions.StartAction(u), nil
}

type Ready struct {
	TypeAction
}

func (a Ready) ToDomain(u *game.User) (fsm.IUserAction, error) {
	return actions.ReadyAction(u), nil
}

func (pc *PlayerConn) UnmarshalAction(msg []byte) (fsm.IUserAction, error) {
	var a IActionInf

	actionType, err := extractTypeAction(msg)
	if err != nil {
		return nil, err
	}

	switch actionType {
	case "answer":
		a = &Answer{}
	case "start":
		a = &Start{}
	case "ready":
		a = &Ready{}
	default:
		return nil, ErrUndefinedAction
	}

	if err := json.Unmarshal(msg, &a); err != nil {
		return nil, err
	}

	return a.ToDomain(pc.user)
}

func extractTypeAction(msg []byte) (string, error) {
	var bt map[string]string
	if err := json.Unmarshal(msg, &bt); err != nil {
		return "", nil
	}
	if t, ex := bt["action"]; !ex {
		return "", ErrTypeDoesNotExist
	} else {
		return t, nil
	}
}
