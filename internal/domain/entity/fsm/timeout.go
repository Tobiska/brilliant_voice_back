package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Timeout struct { //todo timeout has no any user
	Action
}

func (t Timeout) String() string {
	return fmt.Sprintf("Timeout %s", t.U)
}

func TimeoutAction(u *game.User) Timeout {
	return Timeout{
		Action{
			U: u,
		},
	}
}
