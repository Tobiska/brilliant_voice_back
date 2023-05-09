package logicTimer

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

type TimerAdapter struct {
	CancelCtx context.Context
	TimerCh   chan game.TimerInfo
}

func (ta *TimerAdapter) Close() {
	close(ta.TimerCh)
}

func (ta *TimerAdapter) Send(ctx context.Context, m game.TimerInfo) error {
	if err := ta.CancelCtx.Err(); err != nil {
		ta.Close()
		return game.ErrTimerChClosed
	}
	select {
	case <-ctx.Done():
	case <-ta.CancelCtx.Done():
		ta.Close()
		return game.ErrTimerChClosed
	case ta.TimerCh <- m:
	}
	return nil
}
