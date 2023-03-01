package conn

import (
	"brillian_voice_back/internal/domain/entity/game"
	"errors"
)

var (
	ErrDone = errors.New("error done")
)

type adapterConn struct {
	UpdateCh chan game.Game
	ErrCh    chan error
}

func (a *adapterConn) UpdateGame(g game.Game) {
	a.UpdateCh <- g
}

func (a *adapterConn) SendError(err error) {
	a.ErrCh <- err
}

func (a *adapterConn) Close() {
	a.ErrCh <- ErrDone
}
