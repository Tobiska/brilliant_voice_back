package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

type RoundRunning struct{}

func (r *RoundRunning) Current() string {
	return "round_running"
}

func (r *RoundRunning) Send(g *game.Game, a fsm.IUserAction) (fsm.IState, error) {
	if an, ok := a.(actions.Answer); ok {
		return r.handleAnswer(g, an)
	}

	if to, ok := a.(actions.Timeout); ok {
		return r.handleTimeout(g, to)
	}

	if lu, ok := a.(actions.LeaveUser); ok {
		return r.handleLeaveUser(g, lu)
	}

	return &RoundRunning{}, nil
}

func (r *RoundRunning) handleAnswer(g *game.Game, a actions.Answer) (fsm.IState, error) {
	if isFinish, err := func() (bool, error) {
		err := g.Users.Answer(a.User().ID, a.User().Answer)
		return g.Users.CheckAnswers(), err
	}(); !isFinish {
		return &RoundRunning{}, err
	} else {
		return &WaitUsers{}, err
	}
}

func (r *RoundRunning) handleTimeout(g *game.Game, a actions.Timeout) (fsm.IState, error) {
	return nil, nil
}

func (r *RoundRunning) handleLeaveUser(g *game.Game, a actions.LeaveUser) (fsm.IState, error) {
	if err := func() error {
		return g.DeleteUser(a.U)
	}(); errors.Is(err, game.ErrOwnerLeave) {
		return &Dead{}, err
	} else {
		return &RoundRunning{}, err
	}
}

func (r *RoundRunning) Wait() {}
