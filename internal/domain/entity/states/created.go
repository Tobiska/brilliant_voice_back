package states

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type Created struct{}

func (c *Created) Current() string {
	return "created"
}

func (c *Created) Run(_ *game.Game) fsm.IState {
	return &WaitStart{}
}
