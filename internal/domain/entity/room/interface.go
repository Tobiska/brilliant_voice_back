package room

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

var (
	ErrQueueIsEmpty = errors.New("queue is empty")
)

type (
	IPriorityQueue interface {
		Push(a fsm.IAction, p int64)
		Size() int
		Pop() (fsm.IAction, error)
	}

	IClientManager interface {
		Join(user game.User)
	}
)
