package conn

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type PlayerConn struct {
	ws *websocket.Conn

	roomDesc game.Descriptor

	actionCh chan fsm.IAction
	updateCh chan game.Game
	errCh    chan error

	ctxClose context.Context
	cancel   func()
}

func (pc *PlayerConn) Adapter() game.IConn {
	return &adapterConn{
		UpdateCh: pc.updateCh,
		ErrCh:    pc.errCh,
	}
}

func NewPlayerConn(ws *websocket.Conn, roomDesc game.Descriptor, actionCh chan fsm.IAction) *PlayerConn {
	ctx, cancel := context.WithCancel(context.Background())
	pc := &PlayerConn{
		ws:       ws,
		roomDesc: roomDesc,
		actionCh: actionCh,
		ctxClose: ctx,
		cancel:   cancel,
		updateCh: make(chan game.Game),
		errCh:    make(chan error),
	}
	pc.readPump()
	pc.writePump()
	pc.errorPump()
	return pc
}

func (pc *PlayerConn) readPump() {
	go func() {
		for {
			_, m, err := pc.ws.ReadMessage() //todo handle msg
			if err != nil {
				log.Error().Err(err)
				return
			}

			a, err := pc.UnmarshalAction(m)
			if err != nil {
				log.Error().
					Err(err).Msg("unmarshal action error")
				_ = pc.WriteError(err)
				continue
			}

			pc.actionCh <- a
		}
	}()
}

func (pc *PlayerConn) writePump() {
	go func() {
		for {
			select {
			case <-pc.ctxClose.Done():
				return
			case s, ok := <-pc.updateCh:
				if !ok {
					return
				}
				inf, err := ToInfState(s)
				if err != nil {
					log.Error().Err(err).Str("room_id", pc.roomDesc.Code).
						Msg("error occur parse game state")
					continue
				}
				if err := pc.ws.WriteJSON(inf); err != nil {
					log.Error().Err(err).Str("room_id", pc.roomDesc.Code).Msg("error occur send state")
					continue
				}
			}
		}
	}()
}

func (pc *PlayerConn) Write(msg []byte) (int, error) {
	if err := pc.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(msg), nil
}

func (pc *PlayerConn) errorPump() {
	go func() {
		for {
			select {
			case <-pc.ctxClose.Done():
				if err := pc.Close(); err != nil {
					log.Error().Err(err).Str("room_id", pc.roomDesc.Code).Msg("close error")
				}
			case err := <-pc.errCh:
				if pc.handleError(err) {
					pc.cancel()
					return
				}
			}
		}
	}()
}

func (pc *PlayerConn) handleError(err error) bool {
	if err := pc.WriteError(err); err != nil {
		log.Error().Err(err).Str("room_id", pc.roomDesc.Code).Msg("error occur send error")
	}

	return errors.Is(err, ErrDone)
}

func (pc *PlayerConn) Close() error {
	close(pc.errCh)
	close(pc.updateCh)
	return pc.ws.Close()
}
