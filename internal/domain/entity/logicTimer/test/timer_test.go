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
	testCtx := context.Background()
	actionCh := make(chan fsm.IUserAction)
	m := logicTimer.NewManager()
	m.Init(testCtx, actionCh)
	if err := m.Adapter().Send(testCtx, game.TimerInfo{
		TimeOutPeriod: 5,
		TickerPeriod:  5,
	}); err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, <-actionCh)
}

func TestLogicTimerStopped(t *testing.T) {
	testCtx := context.Background()
	actionCh := make(chan fsm.IUserAction)
	m := logicTimer.NewManager()
	m.Init(testCtx, actionCh)
	if err := m.Adapter().Send(context.Background(), game.TimerInfo{
		TimeOutPeriod: 10000000,
		TickerPeriod:  5,
	}); err != nil {
		t.Fatal(err)
	}
	_ = m.Adapter().Send(context.Background(), game.TimerInfo{
		StopFlag: true,
	})
	_ = m.Adapter().Send(context.Background(), game.TimerInfo{
		StopFlag: true,
	})
}

func TestManagerCancel(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	actionCh := make(chan fsm.IUserAction)
	m := logicTimer.NewManager()
	m.Init(cancelCtx, actionCh)
	if err := m.Adapter().Send(context.Background(), game.TimerInfo{
		TimeOutPeriod: 10000000,
		TickerPeriod:  5,
	}); err != nil {
		t.Fatal(err)
	}
	cancel()
	assert.Error(t, m.Adapter().Send(context.Background(), game.TimerInfo{
		StopFlag: true,
	}))
}
