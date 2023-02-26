package gameManager

import (
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/states"
	"brillian_voice_back/internal/domain/entity/user"
)

type GameManager struct {
	numCurRound int

	desc fsm.Descriptor

	fsm *fsm.Fsm
}

func NewManager(code, ownerId string, prop properties.Properties) *GameManager {
	desc := fsm.Descriptor{
		Code:       code,
		IsFully:    false,
		Properties: prop,
	}
	return &GameManager{
		numCurRound: 0,
		desc:        desc,
		fsm: fsm.InitFsm(&states.Created{}, &fsm.Game{
			Descriptor: desc,
			Users:      make(map[string]*user.User, 0),
			OwnerId:    ownerId,
		}),
	}
}

func (m *GameManager) Do(a fsm.IAction) (err error) {
	return m.fsm.SendAction(a)
}

func (m *GameManager) GameDesc() fsm.Descriptor {
	return m.desc
}
