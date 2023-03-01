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
	Send(game *game.Game, a IAction) IState
	Wait()
}

type IActiveState interface {
	IState
	Run(game *game.Game) IState
}

type IAction interface {
	fmt.Stringer
}
