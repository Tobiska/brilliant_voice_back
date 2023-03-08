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

	start chan game.TimerInfo
	stop  chan struct{}

	actionCh chan fsm.IUserAction

	cancelCtx context.Context
}

func NewManager(cancelCtx context.Context, actionCh chan fsm.IUserAction) *Manager {
	m := &Manager{
		start:     make(chan game.TimerInfo),
		stop:      make(chan struct{}),
		cancelCtx: cancelCtx,
		actionCh:  actionCh,
	}
	go m.run()
	return m
}

func (m *Manager) StartCh() chan game.TimerInfo {
	return m.start
}

func (m *Manager) StopCh() chan struct{} {
	return m.stop
}

func (m *Manager) Close() {
	close(m.stop)
	close(m.start)
}

func (m *Manager) run() {
	defer m.Close()
	for {
		if err := m.cancelCtx.Err(); err != nil {
			log.Info().Msg("Timer manager closed")
			return
		}
		select {
		case <-m.cancelCtx.Done():
			log.Info().Msg("Timer manager closed")
			return
		case ti := <-m.start:
			m.handleTimerStart(ti)
		case <-m.stop:
			m.handleTimerStop()
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
	startTime time.Time

	actionCh chan fsm.IUserAction

	cancelCtx context.Context
}

func (lt *LogicTimer) runHandleTimeout() {
	defer lt.timer.Stop()
	select {
	case <-lt.cancelCtx.Done():
		log.Info().Msg("Timer stopped")
		return
	case <-lt.timer.C:
		log.Info().Msg("Timer handle timeout")
		lt.actionCh <- actions.TimeoutAction(game.NewUser("***", "system", &conn.MockConn{})) //todo убрать костыль с системным юзером(*)
		return
	}
}

func (lt *LogicTimer) stop() {
	lt.timer.Stop()
	log.Info().Str("time", time.Now().Sub(lt.startTime).String()).Msg("Current timer stopped")
}
