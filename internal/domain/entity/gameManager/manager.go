package gameManager

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/states"
	"brillian_voice_back/internal/domain/entity/user"
	"github.com/rs/zerolog/log"
)

type GameManager struct {
	numCurRound int

	fsm *fsm.Fsm
}

func NewManager(code, ownerId string, prop properties.Properties) *GameManager {
	return &GameManager{
		numCurRound: 0, //todo init value
		fsm: fsm.InitFsm(&states.Created{}, &fsm.Game{
			Descriptor: fsm.Descriptor{
				Code:       code,
				IsFully:    false,
				Properties: prop,
			},
			Users:   make(map[string]*user.User, 0),
			OwnerId: ownerId,
		}),
	}
}

func (m *GameManager) FinishRound() error {
	return nil
}

func (m *GameManager) HandleLeave(id string) error {
	if m.state.OwnerId == id {
		m.state.status = Dead
		return ErrRoomIsDead
	}
	if err := m.leaveFromRoom(id); err != nil {
		log.Err(err).
			Str("user_id", id).
			Msg("error when trying to leave the room")
	} else {
		log.Info().
			Str("user_id", id).Str("room_id", m.state.Code).
			Msg("user successfully left the room")
	}
	return nil
}

func (m *GameManager) CloseAll() {
	for _, u := range m.state.Users {
		warn := u.DeleteAndClose()
		log.Warn().
			Err(warn).
			Str("user_id", u.ID).
			Msg("occur error when close")
	}
}

func (m *GameManager) AddUser(u *user.User) error {
	//if err := m.check(u); err != nil {
	//	return err
	//}
	//if _, ex := m.state.Users[u.ID]; ex {
	//	return ErrUserAlreadyExist
	//} else {
	//	m.state.Users[u.ID] = u
	//}
	return nil
}

func (m *GameManager) check(user *user.User) error {
	////todo cond if state have status is fully
	//if m.state.Status() == Dead {
	//	return ErrRoomIsDead
	//}
	//if _, err := m.state.GetOwner(); err != nil && user.ID != m.state.OwnerId {
	//	return err
	//}
	return nil
}

func (m *GameManager) leaveFromRoom(id string) error {
	//u, ok := m.state.Users[id]
	//if !ok {
	//	return ErrUserDoesNotExist
	//}
	//delete(m.state.Users, u.ID)
	return nil
}

func (m *GameManager) Do(a actions.IAction) (err error) {
	return m.fsm.SendAction(a)
}
