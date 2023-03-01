package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Answer struct {
	U    *game.User
	Text string
}

func (a Answer) String() string {
	return fmt.Sprintf("answer text: %s game: %s", a.Text, a.U)
}
