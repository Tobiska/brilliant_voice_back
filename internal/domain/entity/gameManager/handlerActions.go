package gameManager

import (
	"brillian_voice_back/internal/domain/entity/actions"
)

func (m *GameManager) HandlePing(a actions.IAction) error {
	p, _ := a.(*actions.Ping)
	if u, ok := m.state.Users[p.UserId]; ok {
		_ = u.Pong()
	}
	return nil
}

func (m *GameManager) HandleAnswer(a actions.IAction) error {
	return nil
}

func (m *GameManager) HandleClose(a actions.IAction) error {
	return nil
}

func (m *GameManager) HandleStart(a actions.IAction) error {
	s, _ := a.(*actions.Start)
	if s.UserId != m.state.OwnerId {
		return ErrUserIsNotOwner
	} else {
		return m.Start()
	}
}

func (m *GameManager) HandleReady(a actions.IAction) error {
	return nil
}
