package gameManager

import "brillian_voice_back/internal/domain/entity/actions"

func (m *GameManager) HandlePing(a actions.IAction) error {
	p, _ := a.(*actions.Ping)
	if u, ok := m.state.Users[p.UserId]; ok {
		_ = u.Pong()
	}
	return nil
}
