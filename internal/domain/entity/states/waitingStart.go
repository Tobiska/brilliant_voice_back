package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/user"
)

type WaitStart struct{}

func (ws *WaitStart) Current() string {
	return "wait start"
}

func (ws *WaitStart) Wait() {}

func (ws *WaitStart) Send(g *fsm.Game, a fsm.IAction) fsm.IState {
	if ad, ok := a.(actions.AddUser); ok {
		return handleAddUser(g, ad.U)
	}
	return &WaitStart{}
}

func handleAddUser(g *fsm.Game, u *user.User) fsm.IState {
	if err := func() error {
		if err := g.AddUser(u); err != nil {
			return err
		}
		if len(g.Users) != g.Properties.CountPlayers {
			return ErrNotEnough
		}
		return nil
	}(); err != nil {
		return &WaitStart{}
	}
	return &Ready{}
}
