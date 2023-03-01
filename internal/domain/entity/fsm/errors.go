package fsm

import "errors"

var (
	ErrRoomIsDead       = errors.New("room is dead")
	ErrUserDoesNotExist = errors.New("game with id doesn't exist")
	ErrUndefinedAction  = errors.New("undefined action")
	ErrEndOfGame        = errors.New("end of the game")
	ErrUserIsNotOwner   = errors.New("game is not owner")
)
