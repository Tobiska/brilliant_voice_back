package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type WaitStart struct{}

func (ws *WaitStart) Current() string {
	return "WAIT_START"
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
		err := g.AddUser(u)
		if g.Users.Len() != g.Properties.CountPlayers {
			return false, err
		}
		return true, nil
	}(); !tr {
		return &WaitStart{}, err
	} else {
		return &Ready{}, err
	}
}
