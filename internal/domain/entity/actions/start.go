package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Start struct {
	Action
}

func (s Start) String() string {
	return fmt.Sprintf("start action")
}

func StartAction(u *game.User) Start {
	return Start{
		Action: Action{
			U: u,
		},
	}
}
