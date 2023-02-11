package conn

import (
	"brillian_voice_back/internal/domain/entity/room"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type PlayerConn struct {
	ws   *websocket.Conn
	room *room.Room
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
				pc.WriteError(err)
				continue
			}

			pc.room.OperationChannel() <- a
		}
	}()
}

func (pc *PlayerConn) Write(msg []byte) (int, error) {
	if err := pc.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}
	return len(msg), nil
}

func (pc *PlayerConn) SendState() error {
	state := pc.room.GetState()
	msg, err := ToInfState(state)
	if err != nil {
		return err
	}
	if err := pc.ws.WriteJSON(msg); err != nil {
		_ = pc.Close()
		return err
	}
	return nil
}

func (pc *PlayerConn) RequestToLeave(id string) {
	pc.room.LeaveChannel() <- id
}

func (pc *PlayerConn) Close() error {
	return pc.ws.Close()
}

func NewPlayerConn(ws *websocket.Conn, r *room.Room) *PlayerConn {
	conn := &PlayerConn{
		ws:   ws,
		room: r,
	}
	conn.receive()
	return conn
}
