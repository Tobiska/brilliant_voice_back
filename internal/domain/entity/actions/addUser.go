package actions

import (
	"brillian_voice_back/internal/domain/entity/user"
	"fmt"
)

type AddUser struct {
	U *user.User
}

func (a AddUser) String() string {
	return fmt.Sprintf("Add user %s", a.U)
}
