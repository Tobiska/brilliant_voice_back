package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type WaitStart struct{}

func (ws *WaitStart) Current() string {
	return "wait start"
}

func (ws *WaitStart) Wait(_ *game.Game) {}

func (ws *WaitStart) Send(g *game.Game, a fsm.IUserAction) (fsm.IState, error) {
	if ad, ok := a.(actions.AddUser); ok {
		return ws.handleAddUser(g, ad.U)
	}
	if l, ok := a.(actions.LeaveUser); ok {
		return handleLeaveUser(g, l)
	}
	return &WaitStart{}, nil
}

func (ws *WaitStart) handleAddUser(g *game.Game, u *game.User) (fsm.IState, error) {
	if tr, err := func() (bool, error) {
		if err := g.AddUser(u); err != nil {
			return false, err
		}
		if g.Users.Len() != g.Properties.CountPlayers {
			return false, nil
		}
		return true, nil
	}(); !tr {
		return &WaitStart{}, err
	} else {
		return &Ready{}, err
	}
}
