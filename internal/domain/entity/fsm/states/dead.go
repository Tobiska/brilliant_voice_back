package states

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

var (
	ErrRoomIsDead = errors.New("room is dead")
)

type Dead struct{}

func (d *Dead) Current() string {
	return "dead"
}

func (d *Dead) Send(_ *game.Game, _ fsm.IUserAction) (fsm.IState, error) {
	return &Dead{}, ErrRoomIsDead
}

func (d *Dead) Wait(_ *game.Game) {}
