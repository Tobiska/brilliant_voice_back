package states

import "brillian_voice_back/internal/domain/entity/fsm"

type Created struct{}

func (c *Created) Current() string {
	return "created"
}

func (c *Created) Run(_ *fsm.Game) fsm.IState {
	return &WaitStart{}
}
