package fsm

import "errors"

var (
	ErrOwnerYetNotJoined = errors.New("owner yet not joined")
	ErrOwnerLeave        = errors.New("owner leave")
	ErrRoomIsDead        = errors.New("room is dead")
	ErrUserDoesNotExist  = errors.New("user with id doesn't exist")
	ErrUndefinedAction   = errors.New("undefined action")
	ErrEndOfGame         = errors.New("end of the game")
	ErrUserIsNotOwner    = errors.New("user is not owner")
)
