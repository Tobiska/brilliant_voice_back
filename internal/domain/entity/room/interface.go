package room

import "brillian_voice_back/internal/domain/entity/user"

type IClientManager interface {
	Join(user user.User)
}

type IHub interface {
}
