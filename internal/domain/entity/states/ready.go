package states

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
)

type Ready struct{}

func (r *Ready) Current() string {
	return "ready"
}

func (r *Ready) Wait() {}

func (r *Ready) Send(g *fsm.Game, a fsm.IAction) fsm.IState {
	if s, ok := a.(actions.Start); ok {
		return handleStart(g, s)
	}
	if l, ok := a.(actions.LeaveUser); ok {
		return handleLeaveUser(g, l)
	}
	return &Ready{}
}

func handleStart(g *fsm.Game, a actions.Start) fsm.IState {
	if err := func() error {
		if a.U.ID != g.OwnerId {
			return ErrStartNotOwner
		}
		return nil
	}; err != nil {
		return &Ready{}
	}
	return &Ready{} //todo RunningRound
}

func handleLeaveUser(g *fsm.Game, a actions.LeaveUser) fsm.IState {
	if err := func() error {
		return g.DeleteUser(a.U)
	}; err != nil {
		return &Dead{}
	}
	return &Ready{} // todo mb wait start
}
