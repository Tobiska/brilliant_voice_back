package conn

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"context"
	"errors"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

type PlayerConn struct {
	ws *websocket.Conn

	roomDesc game.Descriptor
	user     *game.User

	actionCh chan fsm.IUserAction
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

func NewPlayerConn(ws *websocket.Conn, actionCh chan fsm.IUserAction) *PlayerConn {
	ctx, cancel := context.WithCancel(context.Background())
	pc := &PlayerConn{
		ws:       ws,
		actionCh: actionCh,
		ctxClose: ctx,
		cancel:   cancel,
		updateCh: make(chan game.Game),
		errCh:    make(chan error),
	}
	return pc
}

func (pc *PlayerConn) Run() {
	go pc.writePump()
	pc.readPump()
}

func (pc *PlayerConn) SetContextInfo(roomDesc game.Descriptor, user *game.User) {
	pc.user = user
	pc.roomDesc = roomDesc
}

func (pc *PlayerConn) readPump() {
	for {
		_, m, err := pc.ws.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("read message error")
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
}

func (pc *PlayerConn) writePump() {
	defer func() {
		if err := pc.close(); err != nil {
			log.Error().Err(err).Str("room_id", pc.roomDesc.Code).
				Msg("error close websocket connection")
		}
	}()
	for {
		if err := pc.ctxClose.Err(); err != nil {
			return
		}

		select {
		case <-pc.ctxClose.Done():
			return
		case s := <-pc.updateCh:
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
		case err := <-pc.errCh:
			if pc.handleError(err) {
				pc.cancel()
			}
		}
	}
}

func (pc *PlayerConn) Write(msg []byte) (int, error) {
	if err := pc.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(msg), nil
}

func (pc *PlayerConn) handleError(err error) bool {
	if err := pc.WriteError(err); err != nil {
		log.Error().Err(err).Str("room_id", pc.roomDesc.Code).Msg("error occur send error")
	}

	return errors.Is(err, game.ErrDone)
}

func (pc *PlayerConn) StopChannel() <-chan struct{} {
	return pc.ctxClose.Done()
}

func (pc *PlayerConn) close() error {
	_ = pc.ws.WriteMessage(websocket.CloseMessage, []byte{})
	pc.cancel()
	return pc.ws.Close()
}
