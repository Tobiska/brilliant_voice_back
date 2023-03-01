package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Start struct {
	U *game.User
}

func (s Start) String() string {
	return fmt.Sprintf("start action")
}
