package room

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/gameManager"
	"brillian_voice_back/internal/domain/entity/properties"
	"github.com/rs/zerolog/log"
)

type Room struct {
	manager  *gameManager.GameManager
	queue    IPriorityQueue
	actionCh chan fsm.IAction
}

func NewRoom(code, ownerId string,
	prop properties.Properties,
	q IPriorityQueue,
) *Room {
	return &Room{
		manager:  gameManager.NewManager(code, ownerId, prop),
		queue:    q,
		actionCh: make(chan fsm.IAction, 0), //mb add buffer
	}
}

func (r *Room) Run() {
	r.pumpReceiver()
	r.pumpHandler()
}

func (r *Room) Desc() game.Descriptor {
	return r.manager.GameDesc()
}

func (r *Room) ActionChannel() chan fsm.IAction {
	return r.actionCh
}

func (r *Room) pumpReceiver() {
	go func() {
		for {
			if a, ok := <-r.actionCh; ok {
				r.queue.Push(a, 1) // todo diff grad priority
				log.Info().
					Str("action", a.String()).
					Int("size", r.queue.Size()).
					Str("room_id", r.manager.GameDesc().Code).
					Msg("pushed action")
			}
		}
	}()
}

func (r *Room) pumpHandler() {
	go func() {
		for {
			a, err := r.queue.Pop()
			if err != nil {
				break
			}
			if err == ErrQueueIsEmpty {
				continue
			}
			if err := r.manager.Do(a); err != nil {
				log.Error().
					Err(err).
					Str("action", a.String()).
					Str("room_id", r.manager.GameDesc().Code).
					Msg("error handle action")
			}
		}
	}()
}

func (r *Room) JoinToRoom(u *game.User) error {
	return r.manager.Do(&actions.AddUser{
		U: u,
	})
}

func (r *Room) LeaveUser(u *game.User) error {
	return r.manager.Do(&actions.LeaveUser{
		U: u,
	})
}
