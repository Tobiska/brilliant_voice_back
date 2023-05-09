package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
)

type Action struct {
	U *game.User
}

func (a Action) User() *game.User {
	return a.U
}
