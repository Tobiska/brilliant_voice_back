package states

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type Finished struct{}

func (f Finished) Run(_ *game.Game) fsm.IState {
	return &Dead{}
}

func (f Finished) Current() string {
	return "FINISHED"
}
