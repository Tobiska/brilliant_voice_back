package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

type RoundRunning struct {
	currentRound *game.Round
}

func (r *RoundRunning) Current() string {
	return "RUNNING"
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
	r.currentRound.CommitAnswer(a.User(), a.Text)

	if !(len(g.Rounds[g.NumberOfRound].Answers) == g.Users.Len()) {
		return r, nil
	}

	if err := r.StopRound(g); err != nil {
		return &Dead{}, err
	}

	if g.NumberOfRound == len(g.Rounds)-1 {
		return &Results{}, nil
	}
	return &WaitUsers{}, nil
}

func (r *RoundRunning) StopRound(g *game.Game) error {
	err := g.StopTimer()
	g.Users.Reset()
	return err
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
		err := r.StopRound(g)
		return &Dead{}, err
	} else {
		return r, err
	}
}

func (r *RoundRunning) Wait(g *game.Game) {
	r.currentRound = g.Rounds[g.NumberOfRound]
}
