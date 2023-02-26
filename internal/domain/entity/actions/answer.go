package actions

import (
	"brillian_voice_back/internal/domain/entity/user"
	"fmt"
)

type Answer struct {
	U    *user.User
	Text string
}

func (a Answer) String() string {
	return fmt.Sprintf("answer text: %s user: %s", a.Text, a.U)
}
