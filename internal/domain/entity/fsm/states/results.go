package states

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type Results struct {
}

func (r *Results) Current() string {
	return "dead"
}

func (c *Results) Run(_ *game.Game) fsm.IState {
	return &Finished{}
}
