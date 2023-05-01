package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type Ready struct{}

func (r *Ready) Current() string {
	return "READY"
}

func (r *Ready) Wait(_ *game.Game) {}

func (r *Ready) Send(g *game.Game, a fsm.IUserAction) (fsm.IState, error) {
	if s, ok := a.(actions.Start); ok {
		return handleStart(g, s)
	}
	if l, ok := a.(actions.LeaveUser); ok {
		return handleLeaveUser(g, l)
	}
	return &Ready{}, nil
}

func handleStart(g *game.Game, a actions.Start) (fsm.IState, error) {
	if err := func() error {
		if a.U.ID != g.OwnerId {
			return ErrStartNotOwner
		}
		return nil
	}(); err != nil {
		return &Ready{}, err
	}
	if err := g.StartTimer(); err != nil {
		return &Dead{}, err
	}
	return &RoundRunning{}, nil
}

func handleLeaveUser(g *game.Game, a actions.LeaveUser) (fsm.IState, error) {
	if err := func() error {
		return g.DeleteUser(a.U)
	}(); err != nil {
		return &Dead{}, err
	}
	return &WaitStart{}, nil
}
