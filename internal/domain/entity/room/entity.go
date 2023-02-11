package room

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/user"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

var (
	ErrUserAlreadyExist  = errors.New("user with id is already exist")
	ErrUserDoesNotExist  = errors.New("user with id doesn't exist")
	ErrRoomIsDead        = errors.New("room is dead")
	ErrOwnerYetNotJoined = errors.New("owner yet not joined")
)

type GameStatus string

var (
	Wait GameStatus = "WAIT"
	Dead GameStatus = "DEAD"
)

type GameState struct {
	Descriptor
	Users   map[string]*user.User
	OwnerId string
	//questions
	status GameStatus
}

func (s *GameState) GetOwner() (*user.User, error) {
	if u, ok := s.Users[s.OwnerId]; ok {
		return u, nil
	} else {
		return nil, ErrOwnerYetNotJoined
	}
}

func (s *GameState) Status() GameStatus {
	return s.status
}

type Descriptor struct {
	Code       string
	IsFully    bool
	Properties properties.Properties
}

type Room struct {
	state GameState
	mu    *sync.RWMutex

	operationCh chan any
	leaveCh     chan string
	errCh       chan error
}

func (r *Room) OperationChannel() chan any {
	return r.operationCh
}

func (r *Room) LeaveChannel() chan string {
	return r.leaveCh
}

func (r *Room) ErrorChannel() chan error {
	return r.errCh
}

func (r *Room) GetState() GameState {
	return r.state
}

func NewRoom(code, ownerId string, prop properties.Properties) *Room {
	return &Room{
		state: GameState{
			Descriptor: Descriptor{
				Code:       code,
				IsFully:    false,
				Properties: prop,
			},
			Users:   make(map[string]*user.User, 0),
			OwnerId: ownerId,
			status:  Wait,
		},
		operationCh: make(chan any, 0),
		leaveCh:     make(chan string),
		errCh:       make(chan error),

		mu: &sync.RWMutex{},
	}
}

func (r *Room) notifyAll() {
	for _, u := range r.state.Users {
		if err := u.Update(); err != nil {
			fmt.Println("Err notify: ", err)
			continue
		}
	}
}

func (r *Room) Desc() Descriptor {
	return Descriptor{
		Code:       r.state.Code,
		IsFully:    r.state.IsFully,
		Properties: r.state.Properties,
	}
}

func (r *Room) Run() chan error {
	go func() {
		defer r.Finish()
		for {
			select {
			case op := <-r.operationCh:
				log.Info().Str("operation", fmt.Sprintf("%s", op)).Msg("received msg")
				r.notifyAll()
			case id := <-r.leaveCh:
				if r.state.OwnerId == id {
					r.state.status = Dead
					return
				}
				if err := r.leaveFromRoom(id); err != nil {
					log.Err(err).
						Str("user_id", id).
						Msg("error when trying to leave the room")
				} else {
					log.Info().
						Str("user_id", id).Str("room_id", r.Desc().Code).
						Msg("user successfully left the room")
				}
				r.notifyAll()
			}
		}
	}()
	return r.errCh
}

func (r *Room) leaveFromRoom(id string) error {
	u, ok := r.state.Users[id]
	if !ok {
		return ErrUserDoesNotExist
	}
	delete(r.state.Users, u.ID)
	return nil
}

func (r *Room) Finish() {
	for _, u := range r.state.Users {
		_ = u.DeleteAndClose()
	}
	log.Warn().
		Str("room_code", r.Desc().Code).
		Msg("room has finished")
}

func (r *Room) JoinToRoom(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.checkRoom(u); err != nil {
		return err
	}
	if _, ex := r.state.Users[u.ID]; ex {
		return ErrUserAlreadyExist
	} else {
		r.state.Users[u.ID] = u
	}
	r.notifyAll()
	return nil
}

func (r *Room) checkRoom(user *user.User) error {
	//todo cond if state have status is fully
	if r.state.Status() == Dead {
		return ErrRoomIsDead
	}
	if _, err := r.state.GetOwner(); err != nil && user.ID != r.GetState().OwnerId {
		return err
	}
	return nil
}
