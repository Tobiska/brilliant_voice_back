package room

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/gameManager"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/user"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

type Room struct {
	manager *gameManager.GameManager
	mu      *sync.RWMutex

	actionCh chan fsm.IAction
	leaveCh  chan string
	errCh    chan error
}

func (r *Room) OperationChannel() chan fsm.IAction {
	return r.actionCh
}

func (r *Room) LeaveChannel() chan string {
	return r.leaveCh
}

func (r *Room) ErrorChannel() chan error {
	return r.errCh
}

func (r *Room) GetState() gameManager.GameState {
	return r.manager.State()
}

func NewRoom(code, ownerId string, prop properties.Properties) *Room {
	return &Room{
		manager:  gameManager.NewManager(code, ownerId, prop),
		actionCh: make(chan fsm.IAction, 0),
		leaveCh:  make(chan string),
		errCh:    make(chan error),

		mu: &sync.RWMutex{},
	}
}

func (r *Room) notifyAll() {
	r.manager.UpdateAll()
	log.Info().
		Str("room_id", r.Desc().Code).
		Msg("updated all player connections")
}

func (r *Room) Desc() gameManager.Descriptor {
	return r.manager.GameDesc()
}

func (r *Room) Run() chan error {
	go func() {
		defer r.Finish()
		for {
			select {
			case <-r.manager.IsRoundFinishCh():
				if err := r.manager.FinishRound(); err != nil {
					log.Error().
						Err(err).
						Msg("error in time trans to next round")
					break
				} //todo add  errCh

			case <-r.manager.TickerUpdateCh():
				r.notifyAll()
			case a := <-r.actionCh:
				log.Info().
					Str("room_id", r.Desc().Code).
					Str("action", fmt.Sprintf("%s", a)).
					Msg("room handle action")
				if err := r.manager.Do(a); err != nil {
					log.Error().Err(err).
						Str("room_id", r.Desc().Code).
						Str("action", fmt.Sprintf("%s", a)).
						Msg("handle action error")
					//todo write error msg to client
					continue
				}
				r.notifyAll()
			case id := <-r.leaveCh:
				if err := r.manager.HandleLeave(id); err == nil {
					r.notifyAll()
				}
			}
		}
	}()
	return r.errCh
}

func (r *Room) Finish() {
	r.manager.CloseAll()
	log.Warn().
		Str("room_code", r.Desc().Code).
		Msg("room has finished")
}

func (r *Room) JoinToRoom(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.manager.AddUser(u); err != nil {
		return err
	}
	r.notifyAll()
	return nil
}
