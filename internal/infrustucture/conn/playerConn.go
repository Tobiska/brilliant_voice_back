package conn

import (
	"brillian_voice_back/internal/domain/entity/room"
	"fmt"
	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	ws   *websocket.Conn
	room *room.Room
}

func (pc *PlayerConn) receive() {
	go func() {
		for {
			t, _, err := pc.ws.ReadMessage() //todo handle msg
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(t)
			pc.room.OperationChannel() <- "Hey Mambo:)"
		}
	}()
}

func (pc *PlayerConn) Send() error {
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
