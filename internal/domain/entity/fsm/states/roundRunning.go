package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

type RoundRunning struct {
	numberOfRound int
}

func (r *RoundRunning) Current() string {
	return "round_running"
}

func (r *RoundRunning) Send(g *game.Game, a fsm.IUserAction) (fsm.IState, error) {
	if an, ok := a.(actions.Answer); ok {
		return r.handleAnswer(g, an)
	}

	if t, ok := a.(actions.Tick); ok {
		return r.handleTick(g, t)
	}

	if to, ok := a.(actions.Timeout); ok {
		return r.handleTimeout(g, to)
	}

	if lu, ok := a.(actions.LeaveUser); ok {
		return r.handleLeaveUser(g, lu)
	}

	return r, nil
}

func (r *RoundRunning) handleAnswer(g *game.Game, a actions.Answer) (fsm.IState, error) {
	if isFinish, err := func() (bool, error) {
		g.Rounds[r.numberOfRound].Answer(a.User().ID, a.Text)
		return len(g.Rounds[r.numberOfRound].Answers) == g.Users.Len(), nil
	}(); !isFinish {
		return r, err
	} else {
		if err := g.StopTimer(); err != nil {
			return &Dead{}, err
		}
		return &WaitUsers{}, err
	}
}

func (r *RoundRunning) handleTimeout(_ *game.Game, _ actions.Timeout) (fsm.IState, error) {
	//todo проверить system user
	return &WaitUsers{}, nil
}

func (r *RoundRunning) handleTick(g *game.Game, t actions.Tick) (fsm.IState, error) {
	g.RestTime = t.RestOfTime
	return r, nil
}

func (r *RoundRunning) handleLeaveUser(g *game.Game, a actions.LeaveUser) (fsm.IState, error) {
	if err := func() error {
		return g.DeleteUser(a.U)
	}(); errors.Is(err, game.ErrOwnerLeave) {
		if err := g.StopTimer(); err != nil {
			return &Dead{}, err
		}
		return &Dead{}, err
	} else {
		return r, err
	}
}

func (r *RoundRunning) Wait(_ *game.Game) {
}
