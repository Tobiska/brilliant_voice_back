package game

import (
	"context"
	"errors"
)

var (
	ErrDone          = errors.New("error done")
	ErrTimerChClosed = errors.New("error start channel closed")
)

type TimerInfo struct {
	TimeOutPeriod int
	TickerPeriod  int

	StopFlag bool
}

type (
	IConn interface {
		UpdateGame(context.Context, Game)
		SendError(context.Context, error)
		Close(context.Context)
	}

	ITimer interface {
		Send(ctx context.Context, m TimerInfo) error
	}
)
