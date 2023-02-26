package conn

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type PlayerConn struct {
	ws *websocket.Conn

	roomDesc fsm.Descriptor

	actionCh      chan fsm.IAction
	updateStateCh chan fsm.Game
	errCh         chan error
}

func (pc *PlayerConn) UpdateCh() chan fsm.Game {
	return pc.updateStateCh
}

func (pc *PlayerConn) ErrCh() chan error {
	return pc.errCh
}

func NewPlayerConn(ws *websocket.Conn, roomDesc fsm.Descriptor, actionCh chan fsm.IAction) *PlayerConn {
	pc := &PlayerConn{
		ws:            ws,
		roomDesc:      roomDesc,
		actionCh:      actionCh,
		updateStateCh: make(chan fsm.Game),
		errCh:         make(chan error),
	}
	pc.receive()
	pc.writeMsg()
	pc.writeErr()
	return pc
}

func (pc *PlayerConn) receive() {
	go func() {
		for {
			_, m, err := pc.ws.ReadMessage() //todo handle msg
			if err != nil {
				log.Error().Err(err)
				break
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

func (pc *PlayerConn) writeMsg() {
	go func() {
		for {
			s := <-pc.updateStateCh
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
	}()
}

func (pc *PlayerConn) writeErr() {
	go func() {
		for {
			err := <-pc.errCh
			if err := pc.WriteError(err); err != nil {
				log.Error().Err(err).Str("room_id", pc.roomDesc.Code).Msg("error occur send error")
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

func (pc *PlayerConn) Close() error {
	return pc.ws.Close()
}
