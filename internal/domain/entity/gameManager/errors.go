package gameManager

import "errors"

var (
	ErrOwnerYetNotJoined = errors.New("owner yet not joined")
	ErrRoomIsDead        = errors.New("room is dead")
	ErrUserDoesNotExist  = errors.New("user with id doesn't exist")
	ErrUserAlreadyExist  = errors.New("user with id is already exist")
	ErrUndefinedAction   = errors.New("undefined action")
)
