package fsm

import (
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

type Fsm struct {
	currentState IState
	state        *game.Game
}

func (f *Fsm) GameUpdateFmt() *game.Game {
	f.state.Descriptor.Status = f.CurrentState().Current()
	return f.state
}

func InitFsm(initState IState, state *game.Game) *Fsm {
	f := &Fsm{
		currentState: initState,
		state:        state,
	}
	f.Transition(f.currentState)
	return f
}

func (f *Fsm) SendAction(a IUserAction) error {
	if is, ok := f.currentState.(IIdleState); ok {
		st, err := is.Send(f.state, a)
		f.Transition(st)
		return err
	} else {
		return errors.New("current state is not idle")
	}
}

func (f *Fsm) Transition(s IState) {
	f.currentState = s
	if as, ok := f.currentState.(IActiveState); ok {
		f.currentState = as.Run(f.state)
	} else if is, ok := f.currentState.(IIdleState); ok {
		is.Wait()
	}
}

func (f *Fsm) CurrentState() IState {
	return f.currentState
}
