package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

type WaitUsers struct {
	numberOfReadyUser int
}

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

	return wu, nil
}

func (wu *WaitUsers) handleLeaveUser(g *game.Game, a actions.LeaveUser) (fsm.IState, error) {
	if err := func() error {
		return g.DeleteUser(a.U)
	}(); errors.Is(err, game.ErrOwnerLeave) {
		return &Dead{}, err
	} else {
		return wu, err
	}
}

func (wu *WaitUsers) handleReady(g *game.Game, r actions.Ready) (fsm.IState, error) {
	if err := g.Users.MarkReady(r.User()); err != nil {
		return wu, err
	}
	wu.numberOfReadyUser++

	if wu.numberOfReadyUser >= g.Properties.CountPlayers {
		if err := g.StartTimer(); err != nil {
			return &Dead{}, err
		}
		g.NumberOfRound++
		return &RoundRunning{}, nil
	}
	return wu, nil
}

func (wu *WaitUsers) Wait(_ *game.Game) {}
