package gameManager

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/states"
	"sync"
)

type GameManager struct {
	numCurRound int

	desc game.Descriptor

	fsm *fsm.Fsm

	mu *sync.Mutex
}

func NewManager(code, ownerId string, prop properties.Properties) *GameManager {
	desc := game.Descriptor{
		Code:       code,
		IsFully:    false,
		Properties: prop,
	}
	return &GameManager{
		numCurRound: 0,
		mu:          &sync.Mutex{},
		desc:        desc,
		fsm: fsm.InitFsm(&states.Created{}, &game.Game{
			Descriptor: desc,
			Users:      make(map[string]*game.User, 0),
			OwnerId:    ownerId,
		}),
	}
}

func (m *GameManager) Do(a fsm.IAction) (err error) {
	m.mu.Lock()
	err = m.fsm.SendAction(a)
	defer m.mu.Unlock()
	return
}

func (m *GameManager) GameDesc() game.Descriptor {
	return m.desc
}
