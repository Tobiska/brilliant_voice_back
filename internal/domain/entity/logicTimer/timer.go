package logicTimer

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/infrustucture/conn"
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type Manager struct {
	timer *LogicTimer

	timerManagerCh chan game.TimerInfo

	actionCh chan fsm.IUserAction

	cancelCtx context.Context
}

func NewManager() *Manager {
	return &Manager{}
}

// Init создаёт адаптер взаимодействующий с окружающим миром
func (m *Manager) Init(ctx context.Context, actionCh chan fsm.IUserAction) {
	m.actionCh = actionCh
	m.cancelCtx = ctx
	m.timerManagerCh = make(chan game.TimerInfo)
	go m.run()
}

func (m *Manager) Adapter() game.ITimer {
	return &TimerAdapter{
		CancelCtx: m.cancelCtx,
		TimerCh:   m.timerManagerCh,
	}
}

func (m *Manager) run() {
	for {
		if err := m.cancelCtx.Err(); err != nil {
			log.Info().Msg("Timer manager closed")
			return
		}
		select {
		case <-m.cancelCtx.Done():
			log.Info().Msg("Timer manager closed")
			return
		case ti := <-m.timerManagerCh:
			if ti.StopFlag {
				m.handleTimerStop()
			} else {
				m.handleTimerStart(ti)
			}
		}
	}
}

func (m *Manager) handleTimerStart(ti game.TimerInfo) {
	if m.timer != nil {
		log.Warn().Msg("Timer is already started, complete the current")
		return
	}
	m.timer = &LogicTimer{
		cancelCtx: m.cancelCtx,
		startTime: time.Now(),
		ticker:    time.NewTicker(time.Duration(ti.TickerPeriod)),
		timer:     time.NewTimer(time.Duration(ti.TimeOutPeriod)),
		actionCh:  m.actionCh,
	}
	go m.timer.runHandleTimeout()
}

func (m *Manager) handleTimerStop() {
	if m.timer == nil {
		log.Warn().Msg("Timer still not running")
		return
	}
	m.timer.stop()
	m.timer = nil
}

type LogicTimer struct {
	timer     *time.Timer
	ticker    *time.Ticker
	startTime time.Time

	actionCh chan fsm.IUserAction

	cancelCtx context.Context
}

func (lt *LogicTimer) runHandleTimeout() {
	defer lt.timer.Stop()
	for {
		if err := lt.cancelCtx.Err(); err != nil {
			return
		}
		select {
		case <-lt.cancelCtx.Done():
			log.Info().Msg("Timer stopped")
			return
		case <-lt.ticker.C:
			log.Info().Msg("Tick handle")
			lt.actionCh <- actions.TickAction(
				game.NewUser("***", "system", &conn.MockConn{}),
				time.Now().Sub(lt.startTime)) //todo убрать костыль с системным юзером(*)
		case <-lt.timer.C:
			log.Info().Msg("Timer handle timeout")
			lt.actionCh <- actions.TimeoutAction(game.NewUser("***", "system", &conn.MockConn{})) //todo убрать костыль с системным юзером(*)
			return
		}
	}
}

func (lt *LogicTimer) stop() {
	lt.timer.Stop()
	lt.ticker.Stop()
	log.Info().Str("time", time.Now().Sub(lt.startTime).String()).Msg("Current timer stopped")
}
