package states

import "brillian_voice_back/internal/domain/entity/fsm"

type RoundRunning struct {
}

func (r *RoundRunning) Current() string {
	return "round_running"
}

func (r *RoundRunning) Send(g *fsm.Game, a fsm.IAction) fsm.IState {
	return nil
}

func (r *RoundRunning) Wait() {}
