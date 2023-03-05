package room

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
)

type (
	IPriorityQueue interface {
		Deq() (*fsm.IUserAction, IPriorityQueue)
		Size() int
		Enq(*fsm.IUserAction) IPriorityQueue
	}

	IClientManager interface {
		Join(user game.User)
	}
)
