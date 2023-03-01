package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type LeaveUser struct {
	U *game.User
}

func (lu LeaveUser) String() string {
	return fmt.Sprintf("leave game %s", lu.U)
}
