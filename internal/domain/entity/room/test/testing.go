package test

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

type Conn struct {
	UpdateCh chan game.Game
	ErrCh    chan error
}

func (a *Conn) UpdateGame(ctx context.Context, g game.Game) {
	select {
	case a.UpdateCh <- g:
	case <-ctx.Done():
	}
}

func (a *Conn) SendError(ctx context.Context, err error) {
	select {
	case a.ErrCh <- err:
	case <-ctx.Done():
	}
}

func (a *Conn) Close(context.Context) {
	a.ErrCh <- game.ErrDone
}
