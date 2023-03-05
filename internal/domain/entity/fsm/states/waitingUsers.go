package states

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type WaitUsers struct{}

func (wu *WaitUsers) Current() string {
	return "wait users"
}

func (wu *WaitUsers) Send(g *game.Game, a fsm.IUserAction) (fsm.IState, error) {
	return nil, nil
}

func (wu *WaitUsers) Wait() {}
