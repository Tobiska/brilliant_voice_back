package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Timeout struct { //todo timeout has no any user
	Action
}

func (t Timeout) String() string {
	return fmt.Sprintf("Add game %s", t.U)
}

func TimeoutAction(u *game.User) AddUser {
	return AddUser{
		Action{
			U: u,
		},
	}
}
