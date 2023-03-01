package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type AddUser struct {
	U *game.User
}

func (a AddUser) String() string {
	return fmt.Sprintf("Add game %s", a.U)
}
