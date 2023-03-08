package test

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/logicTimer"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogicTimerTimeout(t *testing.T) {
	actionCh := make(chan fsm.IUserAction)
	m := logicTimer.NewManager(context.Background(), actionCh)
	m.StartCh() <- game.TimerInfo{
		TimeOutPeriod: 5,
		TickerPeriod:  5,
	}
	assert.NotNil(t, <-actionCh)
}

func TestLogicTimerStopped(t *testing.T) {
	actionCh := make(chan fsm.IUserAction)
	m := logicTimer.NewManager(context.Background(), actionCh)
	m.StartCh() <- game.TimerInfo{
		TimeOutPeriod: 10000000,
		TickerPeriod:  5,
	}
	m.StopCh() <- struct{}{}
	m.StopCh() <- struct{}{}
}

func TestManagerCancel(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	actionCh := make(chan fsm.IUserAction)
	m := logicTimer.NewManager(cancelCtx, actionCh)
	m.StartCh() <- game.TimerInfo{
		TimeOutPeriod: 10000000,
		TickerPeriod:  5,
	}
	cancel()
	_, ok := <-m.StopCh()
	assert.False(t, ok)
}
