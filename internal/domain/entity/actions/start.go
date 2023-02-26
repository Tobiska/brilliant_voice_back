package actions

import (
	"brillian_voice_back/internal/domain/entity/user"
	"fmt"
)

type Start struct {
	U *user.User
}

func (s *Start) String() string {
	return fmt.Sprintf("start action")
}
