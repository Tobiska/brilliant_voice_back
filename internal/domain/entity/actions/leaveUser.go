package actions

import (
	"brillian_voice_back/internal/domain/entity/user"
	"fmt"
)

type LeaveUser struct {
	U *user.User
}

func (lu *LeaveUser) String() string {
	return fmt.Sprintf("leave user %s", lu.U)
}
