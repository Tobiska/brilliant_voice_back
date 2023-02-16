package gameManager

import (
	"brillian_voice_back/internal/domain/entity/fsm"
)

func (m *GameManager) HandlePing(a fsm.IAction) error {
	p, _ := a.(*fsm.Ping)
	if u, ok := m.state.Users[p.UserId]; ok {
		_ = u.Pong()
	}
	return nil
}

func (m *GameManager) HandleAnswer(a fsm.IAction) error {
	return nil
}

func (m *GameManager) HandleClose(a fsm.IAction) error {
	return nil
}

func (m *GameManager) HandleStart(a fsm.IAction) error {
	s, _ := a.(*fsm.Start)

	return nil
}

func (m *GameManager) HandleReady(a fsm.IAction) error {
	return nil
}
