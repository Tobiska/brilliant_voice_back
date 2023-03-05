package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type LeaveUser struct {
	Action
}

func (lu LeaveUser) String() string {
	return fmt.Sprintf("leave game %s", lu.U)
}

func LeaveUserAction(u *game.User) LeaveUser {
	return LeaveUser{
		Action{
			U: u,
		},
	}
}
