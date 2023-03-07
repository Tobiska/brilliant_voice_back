package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

type WaitUsers struct{}

func (wu *WaitUsers) Current() string {
	return "wait users"
}

func (wu *WaitUsers) Send(g *game.Game, a fsm.IUserAction) (fsm.IState, error) {
	if lu, ok := a.(actions.LeaveUser); ok {
		return wu.handleLeaveUser(g, lu)
	}

	if r, ok := a.(actions.Ready); ok {
		return wu.handleReady(g, r)
	}

	return &WaitUsers{}, nil
}

func (wu *WaitUsers) handleLeaveUser(g *game.Game, a actions.LeaveUser) (fsm.IState, error) {
	if err := func() error {
		return g.DeleteUser(a.U)
	}(); errors.Is(err, game.ErrOwnerLeave) {
		return &Dead{}, err
	} else {
		return &RoundRunning{}, err
	}
}

func (wu *WaitUsers) handleReady(g *game.Game, a actions.Ready) (fsm.IState, error) {
	return nil, nil
}

func (wu *WaitUsers) Wait() {}
