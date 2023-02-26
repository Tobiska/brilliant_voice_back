package fsm

import (
	"fmt"
)

type IState interface {
	Current() string
}

type IIdleState interface {
	IState
	Send(game *Game, a IAction) IState
	Wait()
}

type IActiveState interface {
	IState
	Run(game *Game) IState
}

type IAction interface {
	fmt.Stringer
}
