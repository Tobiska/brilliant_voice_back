package fsm

import (
	"errors"
)

type Fsm struct {
	currentState IState
	state        *Game
}

func InitFsm(initState IState, state *Game) *Fsm {
	return &Fsm{
		currentState: initState,
		state:        state,
	}
}

func (f *Fsm) SendAction(a IAction) error {
	if is, ok := f.currentState.(IIdleState); ok {
		f.Transition(is.Send(f.state, a))
		return nil
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
