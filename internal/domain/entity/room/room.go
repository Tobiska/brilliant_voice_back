package room

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/gameManager"
	"brillian_voice_back/internal/domain/entity/logicTimer"
	"brillian_voice_back/internal/domain/entity/properties"
	"context"
	"github.com/rs/zerolog/log"
)

const (
	BufferSize = 100
)

type Room struct {
	manager *gameManager.GameManager

	userActionCh   chan fsm.IUserAction
	systemActionCh chan fsm.IUserAction

	timerManager *logicTimer.Manager

	cancelCtx context.Context
	cancel    func()
}

func NewRoom(code, ownerId string,
	prop properties.Properties,
	rounds []*game.Round,
	timerManager *logicTimer.Manager,
) *Room {
	ctx, cancel := context.WithCancel(context.TODO())
	actionCh := make(chan fsm.IUserAction, BufferSize)
	systemCh := make(chan fsm.IUserAction, BufferSize)
	timerManager.Init(ctx, systemCh)

	return &Room{
		cancelCtx:      ctx,
		cancel:         cancel,
		manager:        gameManager.NewManager(code, ownerId, prop, rounds, timerManager.Adapter()),
		timerManager:   timerManager,
		userActionCh:   actionCh, //mb add buffer
		systemActionCh: systemCh,
	}
}

func (r *Room) Run() {
	go r.pumpReceiver()
}

func (r *Room) Desc() game.Descriptor {
	return r.manager.Game().Descriptor
}

func (r *Room) ActionChannel() chan fsm.IUserAction {
	return r.userActionCh
}

func (r *Room) pumpReceiver() {
	for {
		if err := r.cancelCtx.Err(); err != nil {
			r.Clear()
			return
		}
		select {
		case <-r.cancelCtx.Done():
			r.Clear()
			return
		case a := <-r.userActionCh:
			log.Info().
				Str("user_action", a.String()).
				Str("room_id", r.manager.Game().Code).
				Msg("handle user action")
			if err := r.manager.DoAsync(a); err != nil {
				log.Error().
					Err(err).
					Str("user_action", a.String()).
					Str("room_id", r.manager.Game().Code).
					Msg("error handle user action")
			}
		case a := <-r.systemActionCh:
			log.Info().
				Str("system_action", a.String()).
				Str("room_id", r.manager.Game().Code).
				Msg("handle system action")
			if err := r.manager.DoAsync(a); err != nil {
				log.Error().
					Err(err).
					Str("system_action", a.String()).
					Str("room_id", r.manager.Game().Code).
					Msg("error handle system action")
			}
		}
	}
}

func (r *Room) Clear() {
	log.Error().
		Str("room_id", r.manager.Game().Code).
		Msg("room receiver stopped")
}
func (r *Room) JoinToRoom(u *game.User) error {
	return r.manager.DoSync(actions.AddUserAction(u))
}

func (r *Room) LeaveUser(u *game.User) error {
	return r.manager.DoSync(actions.LeaveUserAction(u))
}
