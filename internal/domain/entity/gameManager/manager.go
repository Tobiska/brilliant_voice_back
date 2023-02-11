package gameManager

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/user"
	"github.com/rs/zerolog/log"
	"time"
)

type GameManager struct {
	numCurRound int
	state       *GameState
	roundTimer  time.Timer
}

func NewManager(code, ownerId string, prop properties.Properties) *GameManager {
	return &GameManager{
		numCurRound: 0, //todo init value
		state: &GameState{
			Descriptor: Descriptor{
				Code:       code,
				IsFully:    false,
				Properties: prop,
			},
			Users:   make(map[string]*user.User, 0),
			OwnerId: ownerId,
			status:  Wait,
		},
	}
}

func (m *GameManager) State() GameState { //todo remove
	return *m.state
}

func (m *GameManager) UpdateAll() {
	for _, u := range m.state.Users {
		if err := u.Update(); err != nil {
			log.Error().
				Str("user_id", u.ID).
				Err(err).Msg("error in time update")
			continue
		}
	}
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
	if err := m.check(u); err != nil {
		return err
	}
	if _, ex := m.state.Users[u.ID]; ex {
		return ErrUserAlreadyExist
	} else {
		m.state.Users[u.ID] = u
	}
	return nil
}

func (m *GameManager) GameDesc() Descriptor {
	return Descriptor{
		Code:       m.state.Code,
		IsFully:    m.state.IsFully,
		Properties: m.state.Properties,
		Status:     string(m.state.status),
	}
}

func (m *GameManager) check(user *user.User) error {
	//todo cond if state have status is fully
	if m.state.Status() == Dead {
		return ErrRoomIsDead
	}
	if _, err := m.state.GetOwner(); err != nil && user.ID != m.state.OwnerId {
		return err
	}
	return nil
}

func (m *GameManager) leaveFromRoom(id string) error {
	u, ok := m.state.Users[id]
	if !ok {
		return ErrUserDoesNotExist
	}
	delete(m.state.Users, u.ID)
	return nil
}

func (m *GameManager) Do(a actions.IAction) error {
	switch a.Type() {
	case actions.AnswerType:
	case actions.ReadyType:
	case actions.CloseType:
	case actions.StartType:
	case actions.PingType:
		if err := m.HandlePing(a); err != nil {
			return err
		}

	default:
		return ErrUndefinedAction
	}
	return nil
}
