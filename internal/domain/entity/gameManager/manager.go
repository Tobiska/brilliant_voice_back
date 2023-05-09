package gameManager

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/fsm/states"
	"brillian_voice_back/internal/domain/entity/game"
	"context"
	"sync"
	"time"
)

type GameManager struct {
	fsm *fsm.Fsm

	mu *sync.Mutex
}

func NewManager(code, ownerId string, prop game.Properties, rounds []*game.Round, adapter game.ITimer) *GameManager {
	return &GameManager{
		mu: &sync.Mutex{},
		fsm: fsm.InitFsm(&states.Created{}, &game.Game{
			Rounds: rounds,
			Timer:  adapter,
			Descriptor: game.Descriptor{
				Code:       code,
				Properties: prop,
			},
			Users:   game.NewUsersContainer(),
			OwnerId: ownerId,
		}),
	}
}

func (m *GameManager) DoAsync(a fsm.IUserAction) error {
	if err := m.do(a); err != nil {
		m.NotifyError(err, a.User())
		return err
	}
	return nil
}

func (m *GameManager) NotifyError(err error, users ...*game.User) {
	for _, u := range users {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Millisecond) //todo avoid ctx leak
		go u.Conn.SendError(ctx, err)
	}
}

func (m *GameManager) UpdateState(g game.Game, users ...*game.User) {
	for _, u := range users {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Millisecond) //todo avoid ctx leak
		u.Conn.UpdateGame(ctx, g)
	}
}

func (m *GameManager) do(a fsm.IUserAction) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	err = m.fsm.SendAction(a)
	m.UpdateState(m.Game(), m.Game().Users.ToSlice()...)
	return
}

func (m *GameManager) DoSync(a fsm.IUserAction) (err error) {
	return m.do(a)
}

func (m *GameManager) Game() game.Game {
	return *m.fsm.GameUpdateFmt()
}
