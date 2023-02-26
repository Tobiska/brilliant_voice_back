package states

import "brillian_voice_back/internal/domain/entity/fsm"

type Dead struct{}

func (d *Dead) Current() string {
	return "dead"
}

func (d *Dead) Send(g *fsm.Game, a fsm.IAction) fsm.IState {
	return &Dead{}
}

func (d *Dead) Wait() {}
