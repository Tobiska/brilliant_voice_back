package states

import "errors"

var (
	ErrNotEnough        = errors.New("count users not enough")
	ErrStartNotOwner    = errors.New("start not owner")
	ErrUserDoesNotExist = errors.New("start not owner")
)
