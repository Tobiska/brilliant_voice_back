package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Ready struct {
	Action
}

func (a Ready) String() string {
	return fmt.Sprintf("Ready %s", a.U)
}

func ReadyAction(u *game.User) Ready {
	return Ready{
		Action{
			U: u,
		},
	}
}
