package fsm

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type IState interface {
	Current() string
}

type IIdleState interface {
	IState
	Send(game *game.Game, a IUserAction) (IState, error)
	Wait()
}

type IActiveState interface {
	IState
	Run(game *game.Game) IState
}

type IAction interface {
	fmt.Stringer
}

type IUserAction interface {
	IAction
	User() *game.User
}
