package game

import (
	"context"
	"errors"
)

var (
	ErrDone = errors.New("error done")
)

type (
	IConn interface {
		UpdateGame(context.Context, Game)
		SendError(context.Context, error)
		Close(context.Context)
	}

	ITimer interface {
		Start(ctx context.Context, m TimerInfo)
		Stop(ctx context.Context)
	}
)
