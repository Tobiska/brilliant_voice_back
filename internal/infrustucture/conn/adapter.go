package conn

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

type adapterConn struct {
	UpdateCh chan game.Game
	ErrCh    chan error
}

func (a *adapterConn) UpdateGame(ctx context.Context, g game.Game) {
	select {
	case a.UpdateCh <- g:
	case <-ctx.Done():
	}
}

func (a *adapterConn) SendError(ctx context.Context, err error) {
	select {
	case a.ErrCh <- err:
	case <-ctx.Done():
	}
}

func (a *adapterConn) Close(ctx context.Context) {
	a.SendError(ctx, game.ErrDone)
}
