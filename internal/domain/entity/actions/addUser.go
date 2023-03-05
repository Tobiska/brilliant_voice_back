package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type AddUser struct {
	Action
}

func (a AddUser) String() string {
	return fmt.Sprintf("Add game %s", a.U)
}

func AddUserAction(u *game.User) AddUser {
	return AddUser{
		Action{
			U: u,
		},
	}
}
